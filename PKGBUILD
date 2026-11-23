# Maintainer: StepanKomis <komis.stepan@gmail.com>
pkgname=fgss
pkgver=0.1.0
pkgrel=1
pkgdesc='Push and run scripts on remote hosts over ssh/scp, using your ssh config'
arch=('x86_64' 'aarch64')
url='https://github.com/StepanKomis/fgss'
license=('Apache-2.0')
depends=('openssh' 'glibc')
makedepends=('go')
backup=('etc/fgss/scripts/health.sh')
source=("$pkgname-$pkgver.tar.gz::$url/archive/refs/tags/v$pkgver.tar.gz")
sha256sums=('SKIP') # run `updpkgsums` once the tag exists

package() {
  cd "fgss-$pkgver"
  export CGO_ENABLED=0
  make install DESTDIR="$pkgdir" PREFIX=/usr VERSION="$pkgver"
  install -Dm644 LICENSE "$pkgdir/usr/share/licenses/$pkgname/LICENSE"
}
