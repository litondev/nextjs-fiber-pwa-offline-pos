if(!self.define){let e,s={};const i=(i,n)=>(i=new URL(i+".js",n).href,s[i]||new Promise((s=>{if("document"in self){const e=document.createElement("script");e.src=i,e.onload=s,document.head.appendChild(e)}else e=i,importScripts(i),s()})).then((()=>{let e=s[i];if(!e)throw new Error(`Module ${i} didn’t register its module`);return e})));self.define=(n,t)=>{const a=e||("document"in self?document.currentScript.src:"")||location.href;if(s[a])return;let c={};const r=e=>i(e,a),u={module:{uri:a},exports:c,require:r};s[a]=Promise.all(n.map((e=>u[e]||r(e)))).then((e=>(t(...e),c)))}}define(["./workbox-1846d813"],(function(e){"use strict";importScripts(),self.skipWaiting(),e.clientsClaim(),e.precacheAndRoute([{url:"/_next/server/middleware-chunks/606.js",revision:"471UWb2iyIHu5IWpTQLSX"},{url:"/_next/server/pages/auth/_middleware.js",revision:"471UWb2iyIHu5IWpTQLSX"},{url:"/_next/server/pages/pos/_middleware.js",revision:"471UWb2iyIHu5IWpTQLSX"},{url:"/_next/static/471UWb2iyIHu5IWpTQLSX/_buildManifest.js",revision:"471UWb2iyIHu5IWpTQLSX"},{url:"/_next/static/471UWb2iyIHu5IWpTQLSX/_middlewareManifest.js",revision:"471UWb2iyIHu5IWpTQLSX"},{url:"/_next/static/471UWb2iyIHu5IWpTQLSX/_ssgManifest.js",revision:"471UWb2iyIHu5IWpTQLSX"},{url:"/_next/static/chunks/526-edb046bf9a41c74d.js",revision:"471UWb2iyIHu5IWpTQLSX"},{url:"/_next/static/chunks/651.cd440d205ca10b23.js",revision:"471UWb2iyIHu5IWpTQLSX"},{url:"/_next/static/chunks/framework-91d7f78b5b4003c8.js",revision:"471UWb2iyIHu5IWpTQLSX"},{url:"/_next/static/chunks/main-76f598cfe439f872.js",revision:"471UWb2iyIHu5IWpTQLSX"},{url:"/_next/static/chunks/pages/_app-b7c14599162b18ee.js",revision:"471UWb2iyIHu5IWpTQLSX"},{url:"/_next/static/chunks/pages/_error-b007970c9960fd0f.js",revision:"471UWb2iyIHu5IWpTQLSX"},{url:"/_next/static/chunks/pages/auth/login-fc281ce156e0a71b.js",revision:"471UWb2iyIHu5IWpTQLSX"},{url:"/_next/static/chunks/pages/auth/register-35cfbd3fabe60f9b.js",revision:"471UWb2iyIHu5IWpTQLSX"},{url:"/_next/static/chunks/pages/index-818a4af2e5d7aa97.js",revision:"471UWb2iyIHu5IWpTQLSX"},{url:"/_next/static/chunks/pages/pos-248042ef4f21ec0c.js",revision:"471UWb2iyIHu5IWpTQLSX"},{url:"/_next/static/chunks/pages/pos/category-c7f556ceab9b424e.js",revision:"471UWb2iyIHu5IWpTQLSX"},{url:"/_next/static/chunks/pages/pos/customer-eace676aa9257024.js",revision:"471UWb2iyIHu5IWpTQLSX"},{url:"/_next/static/chunks/pages/pos/order-aa3cc383ef7ed3d8.js",revision:"471UWb2iyIHu5IWpTQLSX"},{url:"/_next/static/chunks/pages/pos/product-7bb263f42ddcc6bc.js",revision:"471UWb2iyIHu5IWpTQLSX"},{url:"/_next/static/chunks/pages/pos/transaction-37fd6c8d7ab913c3.js",revision:"471UWb2iyIHu5IWpTQLSX"},{url:"/_next/static/chunks/pages/pos/user-8f09d472f7190b35.js",revision:"471UWb2iyIHu5IWpTQLSX"},{url:"/_next/static/chunks/polyfills-5cd94c89d3acac5f.js",revision:"471UWb2iyIHu5IWpTQLSX"},{url:"/_next/static/chunks/webpack-ad0ec6cb5c39d382.js",revision:"471UWb2iyIHu5IWpTQLSX"},{url:"/_next/static/chunks/webpack-middleware-ad0ec6cb5c39d382.js",revision:"471UWb2iyIHu5IWpTQLSX"},{url:"/_next/static/css/091fe7309fd8b707.css",revision:"471UWb2iyIHu5IWpTQLSX"},{url:"/favicon.ico",revision:"c30c7d42707a47a3f4591831641e50dc"},{url:"/vercel.svg",revision:"4b4f1876502eb6721764637fe5c41702"}],{ignoreURLParametersMatching:[]}),e.cleanupOutdatedCaches(),e.registerRoute("/",new e.NetworkFirst({cacheName:"start-url",plugins:[{cacheWillUpdate:async({request:e,response:s,event:i,state:n})=>s&&"opaqueredirect"===s.type?new Response(s.body,{status:200,statusText:"OK",headers:s.headers}):s}]}),"GET"),e.registerRoute(/^https:\/\/fonts\.(?:gstatic)\.com\/.*/i,new e.CacheFirst({cacheName:"google-fonts-webfonts",plugins:[new e.ExpirationPlugin({maxEntries:4,maxAgeSeconds:31536e3})]}),"GET"),e.registerRoute(/^https:\/\/fonts\.(?:googleapis)\.com\/.*/i,new e.StaleWhileRevalidate({cacheName:"google-fonts-stylesheets",plugins:[new e.ExpirationPlugin({maxEntries:4,maxAgeSeconds:604800})]}),"GET"),e.registerRoute(/\.(?:eot|otf|ttc|ttf|woff|woff2|font.css)$/i,new e.StaleWhileRevalidate({cacheName:"static-font-assets",plugins:[new e.ExpirationPlugin({maxEntries:4,maxAgeSeconds:604800})]}),"GET"),e.registerRoute(/\.(?:jpg|jpeg|gif|png|svg|ico|webp)$/i,new e.StaleWhileRevalidate({cacheName:"static-image-assets",plugins:[new e.ExpirationPlugin({maxEntries:64,maxAgeSeconds:86400})]}),"GET"),e.registerRoute(/\/_next\/image\?url=.+$/i,new e.StaleWhileRevalidate({cacheName:"next-image",plugins:[new e.ExpirationPlugin({maxEntries:64,maxAgeSeconds:86400})]}),"GET"),e.registerRoute(/\.(?:mp3|wav|ogg)$/i,new e.CacheFirst({cacheName:"static-audio-assets",plugins:[new e.RangeRequestsPlugin,new e.ExpirationPlugin({maxEntries:32,maxAgeSeconds:86400})]}),"GET"),e.registerRoute(/\.(?:mp4)$/i,new e.CacheFirst({cacheName:"static-video-assets",plugins:[new e.RangeRequestsPlugin,new e.ExpirationPlugin({maxEntries:32,maxAgeSeconds:86400})]}),"GET"),e.registerRoute(/\.(?:js)$/i,new e.StaleWhileRevalidate({cacheName:"static-js-assets",plugins:[new e.ExpirationPlugin({maxEntries:32,maxAgeSeconds:86400})]}),"GET"),e.registerRoute(/\.(?:css|less)$/i,new e.StaleWhileRevalidate({cacheName:"static-style-assets",plugins:[new e.ExpirationPlugin({maxEntries:32,maxAgeSeconds:86400})]}),"GET"),e.registerRoute(/\/_next\/data\/.+\/.+\.json$/i,new e.StaleWhileRevalidate({cacheName:"next-data",plugins:[new e.ExpirationPlugin({maxEntries:32,maxAgeSeconds:86400})]}),"GET"),e.registerRoute(/\.(?:json|xml|csv)$/i,new e.NetworkFirst({cacheName:"static-data-assets",plugins:[new e.ExpirationPlugin({maxEntries:32,maxAgeSeconds:86400})]}),"GET"),e.registerRoute((({url:e})=>{if(!(self.origin===e.origin))return!1;const s=e.pathname;return!s.startsWith("/api/auth/")&&!!s.startsWith("/api/")}),new e.NetworkFirst({cacheName:"apis",networkTimeoutSeconds:10,plugins:[new e.ExpirationPlugin({maxEntries:16,maxAgeSeconds:86400})]}),"GET"),e.registerRoute((({url:e})=>{if(!(self.origin===e.origin))return!1;return!e.pathname.startsWith("/api/")}),new e.NetworkFirst({cacheName:"others",networkTimeoutSeconds:10,plugins:[new e.ExpirationPlugin({maxEntries:32,maxAgeSeconds:86400})]}),"GET"),e.registerRoute((({url:e})=>!(self.origin===e.origin)),new e.NetworkFirst({cacheName:"cross-origin",networkTimeoutSeconds:10,plugins:[new e.ExpirationPlugin({maxEntries:32,maxAgeSeconds:3600})]}),"GET")}));
