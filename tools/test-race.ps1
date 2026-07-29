$ErrorActionPreference = "Stop"

$toolsDir = Split-Path -Parent $MyInvocation.MyCommand.Path
$importLibrary = Join-Path $toolsDir "libkernelbase.a"
$testBinary = Join-Path $toolsDir "go-dtls-race.test.exe"

zig dlltool -d (Join-Path $toolsDir "kernelbase.def") -l $importLibrary -m i386:x86-64
if ($LASTEXITCODE -ne 0) {
    throw "zig dlltool failed"
}

$env:CGO_ENABLED = "1"
$env:CC = "zig cc"
$env:CXX = "zig c++"
$linkerFlags = "-extldflags `"-L$toolsDir -lkernelbase`""

try {
    go test -race -c -ldflags $linkerFlags -o $testBinary .
    if ($LASTEXITCODE -ne 0) {
        exit $LASTEXITCODE
    }

    # Go's Windows TSan runtime reserves a fixed shadow range. Clear the PE
    # ASLR flags on this disposable test executable so that range is usable.
    $image = [System.IO.File]::ReadAllBytes($testBinary)
    $peOffset = [BitConverter]::ToInt32($image, 0x3c)
    $characteristicsOffset = $peOffset + 24 + 0x46
    $characteristics = [BitConverter]::ToUInt16($image, $characteristicsOffset)
    $characteristics = $characteristics -band (-bnot 0x60)
    $encoded = [BitConverter]::GetBytes([uint16]$characteristics)
    [Array]::Copy($encoded, 0, $image, $characteristicsOffset, 2)
    [System.IO.File]::WriteAllBytes($testBinary, $image)

    $testArguments = @("-test.count=1") + $args
    & $testBinary $testArguments
    if ($LASTEXITCODE -ne 0) {
        exit $LASTEXITCODE
    }
} finally {
    Remove-Item -LiteralPath $testBinary, $importLibrary -Force -ErrorAction SilentlyContinue
}
