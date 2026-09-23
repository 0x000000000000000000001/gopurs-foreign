import { mkdtempSync, mkdirSync, readFileSync, writeFileSync, rmSync } from "node:fs";
import { tmpdir } from "node:os";
import { dirname, resolve, join } from "node:path";
import { fileURLToPath } from "node:url";
import { spawnSync } from "node:child_process";

const root = resolve(dirname(fileURLToPath(import.meta.url)), "..");
const work = mkdtempSync(join(tmpdir(), "gopurs-foreign-function-data-"));
try {
  mkdirSync(join(work, "output/gopurs_runtime"), { recursive: true });
  mkdirSync(join(work, "ffi"));
  writeFileSync(join(work, "go.mod"), "module gopurs\n\ngo 1.27.0\n");
  writeFileSync(join(work, "output/gopurs_runtime/runtime.go"), readFileSync(resolve(root, "../gopurs/runtime/runtime.go")));
  writeFileSync(join(work, "ffi/ffi.go"), readFileSync(join(root, "src/Foreign.go")));
  writeFileSync(join(work, "ffi/ffi_test.go"), readFileSync(join(root, "test/function-data_test.go")));
  const result = spawnSync("go", ["test", "-race", "-count=1", "./ffi"], {
    cwd: work, stdio: "inherit", env: { ...process.env, GOWORK: "off" },
  });
  if (result.error) throw result.error;
  if (result.status !== 0) process.exitCode = result.status ?? 1;
} finally {
  rmSync(work, { recursive: true, force: true });
}
