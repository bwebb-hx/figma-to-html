# Figma -> HTML Design System

This repo was made for researching a design system to automatically generate HTML/CSS from Figma designs, using AI.
It aims to record the different methods discovered, but ultimately to define a streamlined system for quickly automating the generation of Figma designs into real code.

## Research Docs TOC

- [Builder.io](./docs/builder-io.md)
  - A figma plugin that uses AI to generate code from Figma designs.
- [Cursor Figma MCP](./docs/cursor-mcp.md)
  - Using a Figma MCP with Cursor to quickly generate Figma designs with Cursor's AI Agent.

## Scripting & Tools TOC

- ["Figmatic" CLI](./cli/figmatic/README.md)
  - A CLI (written in Go) that interacts with Claude Code and automates various figma code-generation tasks.
- [Scripts Directory](./scripts/)
  - Directory with bash scripts for automating code creation with Claude Code. Mainly used during research, but may be useful reference for those wanting to make bash scripts with Claude Code.