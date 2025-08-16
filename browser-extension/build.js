import esbuild from "esbuild";
import fs from "fs-extra";
import path from "path";
import { watch } from "fs";

const isWatch = process.argv.includes("--watch");

const srcDir = path.join(process.cwd(), "src");
const outDir = path.join(process.cwd(), "dist");

const copyStaticFiles = () => {
  fs.copySync(
    path.join(srcDir, "manifest.json"),
    path.join(outDir, "manifest.json")
  );

  fs.copySync(path.join(srcDir, "popup.html"), path.join(outDir, "popup.html"));
  fs.copySync(path.join(srcDir, "popup.css"), path.join(outDir, "popup.css"));
};

copyStaticFiles();

(async () => {
  const buildOptions = {
    entryPoints: [
      path.join(srcDir, "content.ts"),
      path.join(srcDir, "background.ts"),
      path.join(srcDir, "popup.ts"),
    ],
    bundle: true,
    outdir: outDir,
    format: "esm",
    sourcemap: false,
  };

  if (isWatch) {
    const context = await esbuild.context(buildOptions);
    await context.watch();
    console.log("Watching for changes...");

    watch(srcDir, (_eventType, filename) => {
      if (
        filename &&
        ["manifest.json", "popup.html", "popup.css"].includes(filename)
      ) {
        copyStaticFiles();
      }
    });
  } else {
    await esbuild.build(buildOptions);
    console.log("Build finished!");
  }
})();
