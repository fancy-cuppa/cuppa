#!/usr/bin/env python3
"""Builds the icon files from the logo, assets/cuppa.png.

    python tools/icons/build_icons.py

Writes:
  - the menu-bar icon font (a font with two glyphs, U+E000 and U+E001, the left
    and right half of the cup, so together they make a square icon one cell tall)
    into apps/cuppa-desktop/frontend/src and apps/cuppa-web/frontend/src;
  - assets/icons: every size as a png, icon.ico, icon.icns, a maskable icon and a
    Linux hicolor tree, all from the logo as drawn (frame included);
  - the desktop app icons (apps/cuppa-desktop/build);
  - the web icons and manifest (apps/cuppa-web/frontend/public), with the
    logo's SVG as the vector favicon.

Needs: fonttools, pillow, numpy (and brotli, for a woff2 font).
"""
import io
import os
import sys

import numpy as np
from fontTools.fontBuilder import FontBuilder
from fontTools.pens.ttGlyphPen import TTGlyphPen
from PIL import Image, ImageDraw

ROOT = os.path.abspath(os.path.join(os.path.dirname(__file__), '..', '..'))
LOGO = os.path.join(ROOT, 'assets', 'cuppa.png')

# The glyph box, in font units (1000 per em).
ADVANCE = 600          # width of one cell
HALF = 600             # width of half the icon (one whole cell)
ICON = 1200            # height (and total width) of the icon: two cells wide
TOP = 1000            # the icon top edge, above the baseline


def blocks(img):
    """The bounding boxes of the logo's bright blocks, as (x0, y0, x1, y1) in pixels."""
    a = np.asarray(img.convert('RGBA')).astype(int)
    mask = (a[:, :, 3] > 200) & (a[:, :, 1] > 140) & (a[:, :, 0] < 140)
    h, w = mask.shape
    seen = np.zeros_like(mask)
    out = []
    ys, xs = np.nonzero(mask)
    for y, x in zip(ys, xs):
        if seen[y, x]:
            continue
        stack = [(y, x)]
        seen[y, x] = True
        x0 = x1 = x
        y0 = y1 = y
        count = 0
        while stack:
            cy, cx = stack.pop()
            count += 1
            x0, x1, y0, y1 = min(x0, cx), max(x1, cx), min(y0, cy), max(y1, cy)
            for ny, nx in ((cy + 1, cx), (cy - 1, cx), (cy, cx + 1), (cy, cx - 1)):
                if 0 <= ny < h and 0 <= nx < w and mask[ny, nx] and not seen[ny, nx]:
                    seen[ny, nx] = True
                    stack.append((ny, nx))
        if count > 200:
            out.append((x0, y0, x1 + 1, y1 + 1))
    return out


def crop_box(rects, margin=0.07):
    """A square around the cup, a little larger than its bounding box."""
    x0 = min(r[0] for r in rects)
    y0 = min(r[1] for r in rects)
    x1 = max(r[2] for r in rects)
    y1 = max(r[3] for r in rects)
    side = max(x1 - x0, y1 - y0) * (1 + 2 * margin)
    cx, cy = (x0 + x1) / 2, (y0 + y1) / 2
    return cx - side / 2, cy - side / 2, side


def build_font(rects, path_ttf, path_woff2):
    left, top, side = crop_box(rects)
    scale = ICON / side

    def to_icon(r):  # pixel box to icon units, y measured down from the icon's top
        return ((r[0] - left) * scale, (r[1] - top) * scale, (r[2] - left) * scale, (r[3] - top) * scale)

    def half(which):
        pen = TTGlyphPen(None)
        for r in rects:
            x0, y0, x1, y1 = to_icon(r)
            if which == 0:      # left half: icon x in [0, HALF], drawn at the right of the cell
                lo, hi, shift = 0, HALF, ADVANCE - HALF
            else:               # right half: icon x in [HALF, ICON], drawn at the left of the cell
                lo, hi, shift = HALF, ICON, -HALF
            a, b = max(x0, lo), min(x1, hi)
            if b - a < 2:
                continue
            gx0, gx1 = round(a + shift), round(b + shift)
            gy1, gy0 = round(TOP - y0), round(TOP - y1)
            pen.moveTo((gx0, gy0))
            pen.lineTo((gx0, gy1))
            pen.lineTo((gx1, gy1))
            pen.lineTo((gx1, gy0))
            pen.closePath()
        return pen.glyph()

    empty = TTGlyphPen(None).glyph()
    fb = FontBuilder(1000, isTTF=True)
    names = ['.notdef', 'uniE000', 'uniE001']
    fb.setupGlyphOrder(names)
    fb.setupCharacterMap({0xE000: 'uniE000', 0xE001: 'uniE001'})
    fb.setupGlyf({'.notdef': empty, 'uniE000': half(0), 'uniE001': half(1)})
    fb.setupHorizontalMetrics({n: (ADVANCE, 0) for n in names})
    fb.setupHorizontalHeader(ascent=TOP, descent=-(ICON - TOP))
    fb.setupNameTable({'familyName': 'Cuppa Icons', 'styleName': 'Regular',
                       'uniqueFontIdentifier': 'CuppaIcons-Regular', 'fullName': 'Cuppa Icons',
                       'psName': 'CuppaIcons-Regular', 'version': 'Version 1.0'})
    fb.setupOS2(sTypoAscender=TOP, sTypoDescender=-(ICON - TOP), usWinAscent=TOP, usWinDescent=ICON - TOP,
                fsType=0, achVendID='CUPA')
    fb.setupPost(isFixedPitch=1)
    fb.save(path_ttf)
    try:
        from fontTools.ttLib import TTFont
        f = TTFont(path_ttf)
        f.flavor = 'woff2'
        f.save(path_woff2)
    except Exception as e:  # brotli missing: fall back to woff (zlib)
        print('woff2 not available (%s); writing woff' % e, file=sys.stderr)
        from fontTools.ttLib import TTFont
        f = TTFont(path_ttf)
        f.flavor = 'woff'
        f.save(path_woff2.replace('.woff2', '.woff'))


def logo_square(img):
    """The logo as a 1024 by 1024 picture (the source is a pixel taller)."""
    return img.convert('RGBA').resize((1024, 1024), Image.LANCZOS)


def sized(full, size):
    return full.resize((size, size), Image.LANCZOS)


def maskable(full, size, ground):
    """The logo on a solid ground with the margin Android's mask needs (the
    safe zone is the middle 80%)."""
    out = Image.new('RGBA', (size, size), ground + (255,))
    inner = round(size * 0.72)
    out.paste(sized(full, inner), ((size - inner) // 2, (size - inner) // 2), sized(full, inner))
    return out


def save_ico(path, full, sizes):
    sized(full, 256).save(path, sizes=[(s, s) for s in sizes])


def save_icns(path, full):
    full.save(path, format='ICNS', sizes=[(s, s) for s in (16, 32, 64, 128, 256, 512, 1024)])


MANIFEST = """{
  "name": "Cuppa",
  "short_name": "Cuppa",
  "description": "A visual designer for Bubble Tea terminal interfaces",
  "start_url": "./",
  "display": "standalone",
  "background_color": "#07130d",
  "theme_color": "#07130d",
  "icons": [
    { "src": "android-chrome-192.png", "sizes": "192x192", "type": "image/png" },
    { "src": "android-chrome-512.png", "sizes": "512x512", "type": "image/png" },
    { "src": "maskable-512.png", "sizes": "512x512", "type": "image/png", "purpose": "maskable" },
    { "src": "favicon.svg", "sizes": "any", "type": "image/svg+xml" }
  ]
}
"""


def main():
    img = Image.open(LOGO)
    rects = blocks(img)
    print('%d blocks found' % len(rects))
    if len(rects) < 10:
        sys.exit('the logo changed: expected the cup to be made of bright blocks')
    full = logo_square(img)
    ground = full.getpixel((200, 200))[:3]  # the dark green inside the frame

    for app in ('cuppa-desktop', 'cuppa-web'):
        src = os.path.join(ROOT, 'apps', app, 'frontend', 'src')
        build_font(rects, os.path.join(src, 'cuppa-icons.ttf'), os.path.join(src, 'cuppa-icons.woff2'))

    # assets/icons: the whole set, to take from.
    out = os.path.join(ROOT, 'assets', 'icons')
    os.makedirs(out, exist_ok=True)
    for size in (16, 24, 32, 48, 64, 96, 128, 180, 192, 256, 384, 512, 1024):
        sized(full, size).save(os.path.join(out, 'icon-%d.png' % size))
    save_ico(os.path.join(out, 'icon.ico'), full, (16, 24, 32, 48, 64, 128, 256))
    save_icns(os.path.join(out, 'icon.icns'), full)
    maskable(full, 512, ground).save(os.path.join(out, 'maskable-512.png'))
    for size in (16, 22, 24, 32, 48, 64, 128, 256, 512):
        d = os.path.join(out, 'hicolor', '%dx%d' % (size, size), 'apps')
        os.makedirs(d, exist_ok=True)
        sized(full, size).save(os.path.join(d, 'cuppa.png'))

    # Desktop app (Wails reads build/appicon.png; the others are used as they are).
    build = os.path.join(ROOT, 'apps', 'cuppa-desktop', 'build')
    sized(full, 1024).save(os.path.join(build, 'appicon.png'))
    save_ico(os.path.join(build, 'windows', 'icon.ico'), full, (16, 24, 32, 48, 64, 128, 256))
    save_icns(os.path.join(build, 'darwin', 'icon.icns'), full)

    # Web.
    pub = os.path.join(ROOT, 'apps', 'cuppa-web', 'frontend', 'public')
    for name in ('favicon-192.png', 'favicon-512.png'):  # names an earlier build used
        old = os.path.join(pub, name)
        if os.path.exists(old):
            os.remove(old)
    with open(os.path.join(ROOT, 'assets', 'cuppa.svg'), 'rb') as f:
        svg = f.read()
    with open(os.path.join(pub, 'favicon.svg'), 'wb') as f:
        f.write(svg)
    for size in (16, 32, 48):
        sized(full, size).save(os.path.join(pub, 'favicon-%d.png' % size))
    save_ico(os.path.join(pub, 'favicon.ico'), full, (16, 32, 48))
    sized(full, 180).save(os.path.join(pub, 'apple-touch-icon.png'))
    sized(full, 192).save(os.path.join(pub, 'android-chrome-192.png'))
    sized(full, 512).save(os.path.join(pub, 'android-chrome-512.png'))
    maskable(full, 512, ground).save(os.path.join(pub, 'maskable-512.png'))
    with open(os.path.join(pub, 'site.webmanifest'), 'w', encoding='utf8', newline='\n') as f:
        f.write(MANIFEST)
    print('done')


if __name__ == '__main__':
    main()
