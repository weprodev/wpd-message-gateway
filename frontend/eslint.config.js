import storybook from "eslint-plugin-storybook"

import js from '@eslint/js'
import globals from 'globals'
import reactHooks from 'eslint-plugin-react-hooks'
import reactRefresh from 'eslint-plugin-react-refresh'
import tseslint from 'typescript-eslint'
import { defineConfig, globalIgnores } from 'eslint/config'

/** Feature slices under src/features/ — add here when scaffolding a new slice. */
const FEATURES = ['auth', 'integrations', 'inbox', 'settings', 'workspaces']

const CROSS_FEATURE =
  'Features are independent bounded contexts — do not import another feature. Compose in core/router only.'

const NO_FEATURE_IMPORTS = ['@/features/*', '**/features/*']

function siblingFeaturePatterns(current) {
  return FEATURES.filter((name) => name !== current).flatMap((other) => [
    `@/features/${other}`,
    `@/features/${other}/*`,
    `**/${other}/**`,
  ])
}

export default defineConfig([
  globalIgnores(['dist', 'storybook-static']),
  {
    files: ['**/*.{ts,tsx}'],
    extends: [
      js.configs.recommended,
      tseslint.configs.recommended,
      reactHooks.configs.flat.recommended,
      reactRefresh.configs.vite,
    ],
    languageOptions: {
      ecmaVersion: 2020,
      globals: globals.browser,
    },
  },
  ...FEATURES.map((name) => ({
    files: [`src/features/${name}/**/*.{ts,tsx}`],
    rules: {
      'no-restricted-imports': [
        'error',
        {
          paths: FEATURES.filter((other) => other !== name).map((other) => ({
            name: `@/features/${other}`,
            message: CROSS_FEATURE,
          })),
          patterns: [{ group: siblingFeaturePatterns(name), message: CROSS_FEATURE }],
        },
      ],
    },
  })),
  {
    files: [
      'src/shared/**/*.{ts,tsx}',
      'src/components/**/*.{ts,tsx}',
      'src/lib/**/*.{ts,tsx}',
    ],
    rules: {
      'no-restricted-imports': [
        'error',
        {
          patterns: [
            {
              group: NO_FEATURE_IMPORTS,
              message: 'Must not import features/. Compose in core/router instead.',
            },
          ],
        },
      ],
    },
  },
  {
    files: ['src/core/**/*.{ts,tsx}'],
    ignores: ['src/core/router/**/*.{ts,tsx}'],
    rules: {
      'no-restricted-imports': [
        'error',
        {
          patterns: [
            {
              group: NO_FEATURE_IMPORTS,
              message: 'Only core/router may import feature modules.',
            },
          ],
        },
      ],
    },
  },
  {
    files: ['src/core/router/**/*.{ts,tsx}'],
    rules: {
      'no-restricted-imports': 'off',
    },
  },
  {
    files: ['src/components/ui/**/*.{tsx,ts}'],
    rules: {
      'react-refresh/only-export-components': 'off',
    },
  },
  ...storybook.configs['flat/recommended'],
])
