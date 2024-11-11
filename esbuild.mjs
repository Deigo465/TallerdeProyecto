const fs = require('fs');
const vuePlugin = require("esbuild-plugin-vue3");
const tailwind = require("tailwindcss");
const autoprefixer = require("autoprefixer");
const postCssPlugin = require("@deanc/esbuild-plugin-postcss");
const esbuild = require("esbuild");

const BUILD_OPTIONS = {
  logLevel: 'info',
  entryPoints: ['pkg/web/views/main.ts'],
  bundle: true,
  minify: false,
  sourcemap: true,
  define: {
    'process.env.NODE_ENV': '"development"',
    'process.env.VERSION': '"1.0.0"'
  },
  outfile: 'public/js/app.js',
  minifyWhitespace: true, //avoid leaking routes in vue sfc, maybe we can fix this
  banner: {
    js: '//comment',
    css: '/*comment*/',
  },
  plugins: [
    vuePlugin({
      postcss: {
        plugins: [tailwind, autoprefixer]
      }
      }
    ),
    postCssPlugin({
      plugins: [
        tailwind,
        autoprefixer
      ],
    processOptions: { syntax: require('postcss-scss') },
    }),
  ],
  metafile: true,
}

async function build() {
  let result = await esbuild.build(
    BUILD_OPTIONS
  ).catch(() => process.exit(1))
  fs.writeFileSync('meta.json', JSON.stringify(result.metafile))
}

module.exports = { BUILD_OPTIONS };

const { sourceMapsEnabled } = require('process');
async function watch() {
  let ctx = await esbuild.context(BUILD_OPTIONS);
  await ctx.watch();
  console.log('Watching...');
}
watch();
// build()