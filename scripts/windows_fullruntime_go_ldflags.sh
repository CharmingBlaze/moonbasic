#!/usr/bin/env bash
# Shared by .github/workflows/release.yml and ci.yml for Windows amd64 full-runtime builds.
#
# Contract (PE link):
#   - Go uses -linkmode external so the last link is MinGW g++ driving ld.
#   - -static-libgcc / -static-libstdc++ pull the corresponding libgcc/libstdc++ archive
#     objects into the PE instead of relying on libgcc_s / libstdc++ MinGW DLLs.
#   - winpthread is forced from the static archive between -Wl,-Bstatic … -Wl,-Bdynamic
#     so we do not depend on libwinpthread-1.dll at load time.
#   - Raylib must NOT be linked via -lraylib (no CGO_LDFLAGS): sources compile in-tree.
#   - Jolt is linked via #cgo -lJolt -ljolt_wrapper -lstdc++; static-libstdc++ satisfies
#     the C++ runtime without shipping libstdc++-6.dll next to the exe.
#
# Version injection: set MOONBASIC_WINDOWS_VERSION, or rely on GITHUB_REF_NAME (tags), else "devel".
#
# shellcheck disable=SC2034
moonbasic_windows_fullruntime_go_ldflags() {
  local ver="${MOONBASIC_WINDOWS_VERSION:-${GITHUB_REF_NAME:-devel}}"
  # Multiple -extldflags are how cmd/link forwards argv to the external linker; keep
  # flag order stable (matches tested release.yml behavior).
  printf '%s' \
    "-s -w -X moonbasic/internal/version.Version=${ver} " \
    "-linkmode external " \
    "-extldflags=-static-libgcc " \
    "-extldflags=-static-libstdc++ " \
    "-extldflags=-Wl,-Bstatic " \
    "-extldflags=-lwinpthread " \
    "-extldflags=-Wl,-Bdynamic"
}
