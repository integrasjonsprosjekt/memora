import type { NextConfig } from 'next';
import { loadEnvConfig } from '@next/env';

const projectDir = process.cwd();
loadEnvConfig(projectDir);

const nextConfig: NextConfig = {
  pageExtensions: ['js', 'jsx', 'ts', 'tsx', 'md'],
  output: 'standalone',
  redirects: async () => {
    return [
      {
        source: '/((?!signin|signup).*)',
        destination: '/signin',
        permanent: false,
        missing: [
          {
            type: 'cookie',
            key: '__session',
          }
        ]
      },
    ];
  }
};

export default nextConfig;
