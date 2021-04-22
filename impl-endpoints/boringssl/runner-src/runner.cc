/* Copyright (c) 2014, Google Inc.
 *
 * Permission to use, copy, modify, and/or distribute this software for any
 * purpose with or without fee is hereby granted, provided that the above
 * copyright notice and this permission notice appear in all copies.
 *
 * THE SOFTWARE IS PROVIDED "AS IS" AND THE AUTHOR DISCLAIMS ALL WARRANTIES
 * WITH REGARD TO THIS SOFTWARE INCLUDING ALL IMPLIED WARRANTIES OF
 * MERCHANTABILITY AND FITNESS. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR ANY
 * SPECIAL, DIRECT, INDIRECT, OR CONSEQUENTIAL DAMAGES OR ANY DAMAGES
 * WHATSOEVER RESULTING FROM LOSS OF USE, DATA OR PROFITS, WHETHER IN AN ACTION
 * OF CONTRACT, NEGLIGENCE OR OTHER TORTIOUS ACTION, ARISING OUT OF OR IN
 * CONNECTION WITH THE USE OR PERFORMANCE OF THIS SOFTWARE. */

#include <string>
#include <vector>

#include <openssl/crypto.h>
#include <openssl/err.h>
#include <openssl/ssl.h>

#include <libgen.h>

#include "internal.h"

static const struct argument kArguments[] = {
    {
        "-testcase",
        kRequiredArgument,
        "Handle the tls-interop-runner delegated credentials testcase.",
    },
    {"-as-client", kBooleanArgument,
     "Handle testcases as a client. By default the testcases are handled "
     "from a server perspective."},
    {
        "",
        kOptionalArgument,
        "",
    },
};

int main(int argc, char **argv) {
  CRYPTO_library_init();

  int starting_arg = 1;
  std::vector<std::string> args;
  for (int i = starting_arg; i < argc; i++) {
    args.push_back(argv[i]);
  }

  std::map<std::string, std::string> args_map;
  if (!ParseKeyValueArguments(&args_map, args, kArguments)) {
    return 1;
  }

  if (args_map.count("-as-client") != 0) {
    return DoClient(args_map["-testcase"]);
  }

  return DoServer(args_map["-testcase"]);
}
