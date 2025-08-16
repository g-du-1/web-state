import { beforeEach, describe, expect, it, vi } from "vitest";
import { isUrlDisallowed } from "./isUrlDisallowed";

describe("isUrlDisallowed", () => {
  beforeEach(() => {
    globalThis.chrome = {
      storage: {
        local: {
          get: vi.fn().mockResolvedValue({ whitelistSites: "reddit,facebook" }),
        },
      },
    } as any;
  });

  it("should return true if the url is disallowed", async () => {
    const result = await isUrlDisallowed("https://www.google.com");

    expect(result).toBe(true);
  });

  it('returns true if the list is empty', async () => {
    globalThis.chrome = {
      storage: {
        local: {
          get: vi.fn().mockResolvedValue({ whitelistSites: "" }),
        },
      },
    } as any;

    const result = await isUrlDisallowed("https://www.google.com");

    expect(result).toBe(true);
  })

  it("should return false if the url is allowed", async () => {
    const result = await isUrlDisallowed("https://www.reddit.com");

    expect(result).toBe(false);
  });
});
