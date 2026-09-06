import type { NextConfig } from "next";

const apiOrigin = (
  process.env.RPMP_API_ORIGIN ?? "http://127.0.0.1:8080"
).replace(/\/$/, "");

const nextConfig: NextConfig = {
  reactCompiler: true,
  async rewrites() {
    return [
      {
        source: "/api/v1/:path*",
        destination: `${apiOrigin}/api/v1/:path*`,
      },
    ];
  },
};

export default nextConfig;
