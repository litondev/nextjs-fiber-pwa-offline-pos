const withPWA = require('next-pwa')

module.exports = {
  reactStrictMode: true,
  swcMinify: false,
}

module.exports = withPWA({
  pwa: {
    dest: 'public',
  }
})