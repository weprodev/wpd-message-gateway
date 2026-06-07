import storybook from "eslint-plugin-storybook"

import js from '@eslint/js'
import globals from 'globals'
import reactHooks from 'eslint-plugin-react-hooks'
import reactRefresh from 'eslint-plugin-react-refresh'
import tseslint from 'typescript-eslint'
import { defineConfig, globalIgnores } from 'eslint/config'

export default defineConfig([globalIgnores(['dist', 'storybook-static']), {
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
}, {
  files: ['src/features/**/*.{ts,tsx}'],
  rules: {
    'no-restricted-imports': ['error', {
      patterns: [{
        group: ['@/features/*'],
        message: 'Features are isolated — use relative imports within the same feature. Route composition lives in core/router.',
      }],
    }],
  },
}, {
  files: ['src/shared/**/*.{ts,tsx}'],
  rules: {
    'no-restricted-imports': ['error', {
      patterns: [{
        group: ['@/features/*'],
        message: 'shared/ must not import features/. Compose pages in core/router instead.',
      }],
    }],
  },
}, {
  files: ['src/core/**/*.{ts,tsx}'],
  rules: {
    'no-restricted-imports': ['error', {
      patterns: [{
        group: ['@/features/*'],
        message: 'Only core/router may import feature modules.',
      }],
    }],
  },
}, {
  files: ['src/core/router/**/*.{ts,tsx}'],
  rules: {
    'no-restricted-imports': 'off',
  },
}, {
  files: ['src/components/ui/**/*.{tsx,ts}'],
  rules: {
    'react-refresh/only-export-components': 'off',
  },
}, ...storybook.configs["flat/recommended"]])
