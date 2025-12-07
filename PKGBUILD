# Maintainer: Your Name <you@example.com>

pkgname=penguinguide-bin
pkgver=0.1.0
pkgrel=1
pkgdesc="Friendly Linux helper that teaches commands while doing tasks"
arch=('x86_64' 'aarch64')
url="https://github.com/scottc943/penguinguide"
license=('MIT')
provides=('penguinguide')
conflicts=('penguinguide')
depends=('glibc')
optdepends=(
  'iproute2: additional routing and interface info output'
  'curl: improved public IP detection and extended speed tests'
  'networkmanager: nicer wifi output through nmcli'
  'wireless_tools: fallback wifi data using iwconfig'
)
makedepends=('tar')

source_x86_64=("${url}/releases/download/v${pkgver}/penguinguide_${pkgver}_linux_amd64.tar.gz")
source_aarch64=("${url}/releases/download/v${pkgver}/penguinguide_${pkgver}_linux_arm64.tar.gz")

sha256sums_x86_64=('SKIP')
sha256sums_aarch64=('SKIP')

package() {
  cd "${srcdir}"

  if [[ "$CARCH" == "x86_64" ]]; then
    tar xf "penguinguide_${pkgver}_linux_amd64.tar.gz"
  else
    tar xf "penguinguide_${pkgver}_linux_arm64.tar.gz"
  fi

  # Install main binary
  install -Dm755 "penguinguide" "${pkgdir}/usr/bin/penguinguide"

  # Install man page if present
  if [[ -f "man/penguinguide.1" ]]; then
    install -Dm644 "man/penguinguide.1" "${pkgdir}/usr/share/man/man1/penguinguide.1"
  fi

  install -Dm644 "LICENSE" "${pkgdir}/usr/share/licenses/${pkgname}/LICENSE"
}

