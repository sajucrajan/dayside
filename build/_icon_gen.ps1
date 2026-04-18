# Generates build/appicon.png and build/windows/icon.ico
# Icon concept: an "eye" glyph (revealing the hidden) on an indigo->purple
# gradient rounded-square, with a faint radar ring around the eye.
# Run from project root:  powershell -ExecutionPolicy Bypass -File build\_icon_gen.ps1

$ErrorActionPreference = 'Stop'

Add-Type -AssemblyName System.Drawing

function New-SpectreIconBitmap {
    param([int]$size)

    $bmp = New-Object System.Drawing.Bitmap $size, $size
    $g = [System.Drawing.Graphics]::FromImage($bmp)
    $g.SmoothingMode = [System.Drawing.Drawing2D.SmoothingMode]::AntiAlias
    $g.PixelOffsetMode = [System.Drawing.Drawing2D.PixelOffsetMode]::HighQuality
    $g.InterpolationMode = [System.Drawing.Drawing2D.InterpolationMode]::HighQualityBicubic

    # Rounded-square background
    $radius = [int]($size * 0.22)
    $path = New-Object System.Drawing.Drawing2D.GraphicsPath
    $path.AddArc(0, 0, $radius*2, $radius*2, 180, 90) | Out-Null
    $path.AddArc($size-$radius*2, 0, $radius*2, $radius*2, 270, 90) | Out-Null
    $path.AddArc($size-$radius*2, $size-$radius*2, $radius*2, $radius*2, 0, 90) | Out-Null
    $path.AddArc(0, $size-$radius*2, $radius*2, $radius*2, 90, 90) | Out-Null
    $path.CloseFigure()

    # Indigo->deep-purple gradient (top-left to bottom-right)
    $pt1 = New-Object System.Drawing.Point 0, 0
    $pt2 = New-Object System.Drawing.Point $size, $size
    $c1 = [System.Drawing.ColorTranslator]::FromHtml('#7c3aed')
    $c2 = [System.Drawing.ColorTranslator]::FromHtml('#1e1b4b')
    $bgBrush = New-Object System.Drawing.Drawing2D.LinearGradientBrush $pt1, $pt2, $c1, $c2
    $g.FillPath($bgBrush, $path)

    $cx = $size / 2
    $cy = $size / 2

    # Outer radar ring (very faint)
    $ringPen = New-Object System.Drawing.Pen ([System.Drawing.Color]::FromArgb(60, 255, 255, 255)), ($size * 0.012)
    $ringR = $size * 0.40
    $g.DrawEllipse($ringPen, ($cx - $ringR), ($cy - $ringR), ($ringR * 2), ($ringR * 2))

    # Inner radar ring (slightly brighter)
    $ring2Pen = New-Object System.Drawing.Pen ([System.Drawing.Color]::FromArgb(110, 255, 255, 255)), ($size * 0.014)
    $ring2R = $size * 0.32
    $g.DrawEllipse($ring2Pen, ($cx - $ring2R), ($cy - $ring2R), ($ring2R * 2), ($ring2R * 2))

    # Eye (white ellipse, wider than tall)
    $eyeW = $size * 0.50
    $eyeH = $size * 0.30
    $eyeBrush = New-Object System.Drawing.SolidBrush ([System.Drawing.Color]::White)
    $g.FillEllipse($eyeBrush, ($cx - $eyeW/2), ($cy - $eyeH/2), $eyeW, $eyeH)

    # Pupil (dark indigo)
    $pupilR = $size * 0.10
    $pupilBrush = New-Object System.Drawing.SolidBrush ([System.Drawing.ColorTranslator]::FromHtml('#1e1b4b'))
    $g.FillEllipse($pupilBrush, ($cx - $pupilR), ($cy - $pupilR), ($pupilR * 2), ($pupilR * 2))

    # Inner highlight pupil dot (purple)
    $coreR = $size * 0.035
    $coreBrush = New-Object System.Drawing.SolidBrush ([System.Drawing.ColorTranslator]::FromHtml('#a78bfa'))
    $g.FillEllipse($coreBrush, ($cx - $coreR), ($cy - $coreR), ($coreR * 2), ($coreR * 2))

    # Small specular shine on pupil
    $shineR = $size * 0.022
    $shineBrush = New-Object System.Drawing.SolidBrush ([System.Drawing.Color]::White)
    $g.FillEllipse($shineBrush, ($cx + $pupilR * 0.3 - $shineR), ($cy - $pupilR * 0.5 - $shineR), ($shineR * 2), ($shineR * 2))

    $g.Dispose()
    return $bmp
}

# Write PNG (1024x1024 - Wails reads this as the source)
$png = New-SpectreIconBitmap -size 1024
$scriptDir = Split-Path -Parent $MyInvocation.MyCommand.Path
$pngPath = Join-Path $scriptDir 'appicon.png'
Write-Host "Writing to: $pngPath"
# If target is locked (e.g. running exe, thumbnail cache), write to a .new
# sibling and leave swap-in to the caller.
try {
    if (Test-Path $pngPath) { Remove-Item $pngPath -Force -ErrorAction Stop }
    $png.Save($pngPath, [System.Drawing.Imaging.ImageFormat]::Png)
    Write-Host "Wrote $pngPath"
} catch {
    $pngPath = Join-Path $scriptDir 'appicon.new.png'
    if (Test-Path $pngPath) { Remove-Item $pngPath -Force }
    $png.Save($pngPath, [System.Drawing.Imaging.ImageFormat]::Png)
    Write-Host "Original locked; wrote $pngPath (move it into place when safe)"
}

# Write ICO (multi-resolution, PNG-compressed entries per ICO spec v6)
$icoPath = Join-Path $scriptDir 'windows\icon.ico'
$icoLocked = $false
try {
    if (Test-Path $icoPath) { Remove-Item $icoPath -Force -ErrorAction Stop }
} catch {
    $icoPath = Join-Path $scriptDir 'windows\icon.new.ico'
    if (Test-Path $icoPath) { Remove-Item $icoPath -Force }
    $icoLocked = $true
}
$sizes = @(16, 32, 48, 64, 128, 256)

$ms = New-Object System.IO.MemoryStream
$bw = New-Object System.IO.BinaryWriter $ms

# ICONDIR (6 bytes): reserved=0, type=1 (icon), count
$bw.Write([uint16]0)
$bw.Write([uint16]1)
$bw.Write([uint16]$sizes.Count)

# Entries: render each size, capture PNG bytes, compute offset table
$pngStreams = @()
$offset = 6 + ($sizes.Count * 16)  # header + all dir entries

foreach ($s in $sizes) {
    $b = New-SpectreIconBitmap -size $s
    $pstream = New-Object System.IO.MemoryStream
    $b.Save($pstream, [System.Drawing.Imaging.ImageFormat]::Png)
    $bytes = $pstream.ToArray()
    $pngStreams += ,$bytes

    $w = if ($s -ge 256) { 0 } else { $s }   # 0 means 256 per ICO spec
    $h = if ($s -ge 256) { 0 } else { $s }

    $bw.Write([byte]$w)             # width
    $bw.Write([byte]$h)             # height
    $bw.Write([byte]0)              # color count (0 for truecolor)
    $bw.Write([byte]0)              # reserved
    $bw.Write([uint16]1)            # color planes
    $bw.Write([uint16]32)           # bits per pixel
    $bw.Write([uint32]$bytes.Length)
    $bw.Write([uint32]$offset)

    $offset += $bytes.Length
    $b.Dispose()
    $pstream.Dispose()
}

foreach ($b in $pngStreams) {
    $bw.Write($b)
}

[System.IO.File]::WriteAllBytes($icoPath, $ms.ToArray())
if ($icoLocked) {
    Write-Host "Original icon.ico locked; wrote $icoPath  ($($sizes -join ', ') px). Move into place when safe."
} else {
    Write-Host "Wrote $icoPath  ($($sizes -join ', ') px)"
}

$bw.Dispose()
$ms.Dispose()
$png.Dispose()
