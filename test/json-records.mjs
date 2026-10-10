import { mkdtempSync, mkdirSync, readFileSync, writeFileSync } from "node:fs";
import { tmpdir } from "node:os";
import { dirname, resolve, join } from "node:path";
import { fileURLToPath } from "node:url";
import { spawnSync } from "node:child_process";

const root = resolve(dirname(fileURLToPath(import.meta.url)), "..");
const work = mkdtempSync(join(tmpdir(), "gopurs-foreign-json-records-"));
console.log(`Workspace: ${work}`);
mkdirSync(join(work, "output/gopurs_runtime"), { recursive: true });
mkdirSync(join(work, "ffi"));
writeFileSync(join(work, "go.mod"), "module gopurs\n\ngo 1.27.0\n");
writeFileSync(join(work, "output/gopurs_runtime/runtime.go"), readFileSync(resolve(root, "../gopurs/runtime/runtime.go")));
writeFileSync(join(work, "ffi/ffi.go"), readFileSync(join(root, "src/Foreign.go")));
writeFileSync(join(work, "ffi/ffi_test.go"), readFileSync(join(root, "test/json-records_test.go")));
const result = spawnSync("go", ["test", "-race", "-count=1", "-v", "./ffi"], {
  cwd: work, encoding: "utf8", env: { ...process.env, GOWORK: "off" }, timeout: 120000,
});
writeFileSync(join(work, "test.log"), (result.stdout ?? "") + (result.stderr ?? ""));
process.stdout.write(result.stdout ?? "");
process.stderr.write(result.stderr ?? "");
if (result.error) throw result.error;
process.exitCode = result.status ?? 1;
