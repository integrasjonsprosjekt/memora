import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  pageExtensions: ['js', 'jsx', 'ts', 'tsx', 'md'],
  output: 'standalone',
}

export default nextConfig;
