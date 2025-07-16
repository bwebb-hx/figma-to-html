# Figmatic: A CLI for Generating Code from Figma Designs

This CLI is for facilitating the creation of code from Figma designs. As of now, it only generates HTML and CSS code, but perhaps in the future the scope can be expanded to other programming languages and front-end frameworks, such as React.

## Requirements

You must have `claude` (Claude Code) installed and findable in your system's PATH variable.

## Setup

Install this CLI globally:

```bash
go install github.com/bwebb-hx/figma-to-html/cli/figmatic
```

### Required! MCP Server Setup

You must have the Figma MCP server registered in Claude Code, under the **USER** scope:

```bash
claude mcp add figma-mcp-1 -s user -- npx -y figma-developer-mcp --figma-api-key=$FIGMA_ACCESS_TOKEN --stdio
```

> For now, we have a hard-coded MCP name of "figma-mcp-1". However, in the future I may add config for customizing this name.

If you are unsure if your MCP server is set up right, try entering the following command:

```bash
claude mcp list
```

If it's correctly set under the user scope, it should show your MCP config regardless of what your current directory is in the terminal.

## Usage

For detailed usage instructions, enter `figmatic -h` in the terminal.

### Generate Figma Designs

Here's a simple usage example of the `gen` command:

```bash
figmatic gen --url "https://www.figma.com/design/AxOg233wgtoO27RFd8QIOv/M55_コーディングテスト用?node-id=1-1062" -t "$FIGMA_ACCESS_TOKEN" --sub-nodes
```

This command does the following:

```
gen:          command for generating code from figma designs
--url:        specify the URL of the figma design
-t:           your Figma access token
--sub-nodes:  tells the tool to generate the sub-nodes (one level below) first, and then combine those results. For improved accuracy.
```

#### Flags

See the help text for info about the various flags; things like `--url` are a bit self-explanatory. But, here is some info on the more nuanced and optional flags:

- `--sub-nodes`: Split the given design into the nodes directly below it, and combine the result of each generated sub-node into one HTML file.
  - Using this flag will potentially multiply the execution time, depending on how many sub-nodes there are under the node you provided.
  - This technique results in **much higher quality** for large and complex designs, and is highly recommended for such designs.
  - If you are only generating a small node without many complex sub-nodes below it, then it may not be useful.

- `--iterations`: Specify a number of times for Claude Code to iterate on a generated HTML page, attempting to improve it and fix any inaccuracies.
  - Using this flag will, of course, potentially multiply the execution time of the script.
  - In my own testing, it seems that 1 or 2 iterations can yield good results, but more than that and each iteration begins to yield less improvement.
  - In some cases, too much iteration seems like it could cause **bad** effects, where the design actually looks a little worse as a result.
  - I'd only recommend using it if the design you are generating is quite complicated. If you are generating using the `--sub-nodes` flag, it may be less useful since each individual node being generated is smaller.