#!/usr/bin/env python3
"""Builds the icon files from the logo, assets/cuppa.png.

    python tools/icons/build_icons.py

Writes:
  - the menu-bar icon font (a font with two glyphs, U+E000 and U+E001, the left
    and right half of the cup, so together they make a square icon one cell tall)
    into apps/cuppa-desktop/frontend/src and apps/cuppa-web/frontend/src;
  - the desktop app icon (apps/cuppa-desktop/build/appicon.png and
    build/windows/icon.ico);
  - the web favicon set (apps/cuppa-web/frontend/public).

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


def rounded_square(img, rects, size, radius=0.2, pad=0.0):
    """The cup on its dark ground, cropped square, rounded, at size x size."""
    # The square inside the logo's black frame, centred on the picture.
    w, h = img.size
    side = round(min(w, h) * 0.86)
    left, top = (w - side) // 2, (h - side) // 2
    crop = img.convert('RGBA').crop((left, top, left + side, top + side))
    inner = round(size * (1 - 2 * pad))
    crop = crop.resize((inner, inner), Image.LANCZOS)
    mask = Image.new('L', (inner * 4, inner * 4), 0)
    ImageDraw.Draw(mask).rounded_rectangle((0, 0, inner * 4 - 1, inner * 4 - 1), radius=round(inner * 4 * radius), fill=255)
    mask = mask.resize((inner, inner), Image.LANCZOS)
    out = Image.new('RGBA', (size, size), (0, 0, 0, 0))
    out.paste(crop, ((size - inner) // 2, (size - inner) // 2), mask)
    return out


def main():
    img = Image.open(LOGO)
    rects = blocks(img)
    print('%d blocks found' % len(rects))
    if len(rects) < 10:
        sys.exit('the logo changed: expected the cup to be made of bright blocks')

    for app in ('cuppa-desktop', 'cuppa-web'):
        src = os.path.join(ROOT, 'apps', app, 'frontend', 'src')
        build_font(rects, os.path.join(src, 'cuppa-icons.ttf'), os.path.join(src, 'cuppa-icons.woff2'))

    # Desktop app icon: 1024 with the margin macOS expects, and a full-bleed .ico for Windows.
    build = os.path.join(ROOT, 'apps', 'cuppa-desktop', 'build')
    rounded_square(img, rects, 1024, radius=0.2, pad=0.1).save(os.path.join(build, 'appicon.png'))
    bleed = rounded_square(img, rects, 256, radius=0.18)
    bleed.save(os.path.join(build, 'windows', 'icon.ico'), sizes=[(s, s) for s in (16, 24, 32, 48, 64, 128, 256)])

    # Web favicon set.
    pub = os.path.join(ROOT, 'apps', 'cuppa-web', 'frontend', 'public')
    for size in (16, 32, 48, 192, 512):
        rounded_square(img, rects, size, radius=0.18).save(os.path.join(pub, 'favicon-%d.png' % size))
    rounded_square(img, rects, 180, radius=0.0).save(os.path.join(pub, 'apple-touch-icon.png'))
    rounded_square(img, rects, 256, radius=0.18).save(
        os.path.join(pub, 'favicon.ico'), sizes=[(16, 16), (32, 32), (48, 48)])
    print('done')


if __name__ == '__main__':
    main()
