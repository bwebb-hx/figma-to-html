# Using Claude Code to Generate HTML from Figma

## Adding an MCP Server to Claude Code

Related docs: https://docs.anthropic.com/en/docs/claude-code/mcp#configure-mcp-servers

Here's the basic command syntax for adding a new mcp server to claude code:

```
claude mcp add <name> -- <command> [args...]
```

> Note: in the documentation, they don't show the "--" in the "Basic Syntax" section, but they do in the example.
> I tried without the "--" first and it errored, so it seems that is probably required.

For our Figma MCP server command, we can enter something like the following:

```
claude mcp add figma-mcp-1 -- npx -y figma-developer-mcp --figma-api-key=FIGMA_ACCESS_TOKEN_HERE --stdio
```

> Make sure to sub in your actual figma access token

See the full documentation for other relevant MCP commands.

## Attempting to use Figma MCP Server with Claude Code

It looks like the MCP server is working, so I gave Claude Code a try by giving it a figma design link.

I've gotten an error once the MCP server started downloading the figma design data:

```
⏺ figma-mcp-1:get_figma_data (MCP)(fileKey: "AxOg233wgtoO27RFd8QIOv", nodeId: "1-1061")
  ⎿  Error: MCP tool "get_figma_data" response (27548 tokens) exceeds maximum allowed tokens (25000). Please use
     pagination, filtering, or limit parameters to reduce the response size.
```

I simply ran the same thing again (I think? maybe claude code adjusted the approach this time?) and it started working now.

It successfully generated some HTML, but it seems to have only created a small amount of the design - mainly the top section and only a tiny bit of it.

I asked it to take another look at the design and try again. 

> Did it fail to download all of the design due to that previous error about "maximum allowed tokens"?

After it finished the second iteration, much more of the design was now there. Awesome! But, still a significant amount of inaccuracies. Still, better than nothing. I wonder if I'll need to do several iterations to slowly get closer and closer to the correct design? Let's try one more iteration.

The third iteration has taken significantly longer (around 300s), and I did notice an error message in the output:

```
⏺ Update(index.html)
  ⎿  Error: Found 4 matches of the string to replace, but replace_all is false. To replace all occurrences, set
     replace_all to true. To replace only one occurrence, please provide more context to uniquely identify the
     instance.
     String:                         <div class="news-item">
                                 <div class="news-date">2025年6月9日</div>
                                 <div class="news-content">
                                     <span class="news-label">NEWS</span>
                                     <div class="news-text">Consolidated Financial Results for the Fiscal Year Ended
     March 31, 2025 (Under Japanese GAAP)</div>
                                 </div>
                             </div>
```

It seemed like Claude Code was able to move on from this and try to resolve the issue. At some point it did finish up, but not clear what the hold up was.
Upon reviewing the 3rd iteration, it did seem to fix a couple things and get a bit closer to the figma design, but I think the improvements have slowed down.

Maybe this tool will just be mainly useful for getting designs started, and then manual intervention will be necessary? Or maybe it's better to build smaller pieces bit by bit, and put it all together?

After this session with 3 iterations of asking claude code to generate HTML for the figma design, it looks like in total we incurred about $1.50 of charges:

```
Total cost:            $1.49
Total duration (API):  21m 25.5s
Total duration (wall): 37m 7.4s
Total code changes:    781 lines added, 57 lines removed
Usage by model:
    claude-3-5-haiku:  68.9k input, 2.5k output, 0 cache read, 0 cache write
       claude-sonnet:  198 input, 23.6k output, 1.2m cache read, 186.9k cache write
```

## Creating a Script to Automate Figma Code Generation

The main benefit of claude code over cursor is that we can potentially use this CLI version of claude to create scripts and automate things. Let's explore the possibilities here.

First, let's think of how scripting or automation would be useful in terms of converting figma designs to code:

1. Given a list of URLs, generate the HTML/CSS for each design.
- this would be especially useful if there were a lot of designs we need to generate, and each design might take some time.
- give the script a list of URLs and let it run (with auto-accept edits turned on). come back later and find all the designs converted into HTML!
- if generating smaller sub-sections of the design is more accurate, we could also divide a page of a figma design into sub-sections, get the URLs for each, have those generated individually, and then put them all together after.

2. An iterative claude code session?
- I observed in my first test that the first iteration was missing a large amount of content and generally inaccurate.
- So, perhaps we could make a script that causes claude to iterate a few times on the same figma design, gradually improving it without needing user intervention.
- We could give claude a prompt each time, perhaps asking it to "look for inaccuracies and fix them". And set a number of iterations to repeat this?

### Giving Permissions to Claude Code Non-interactive Mode

To use Claude Code in a script, you can simply call it and give it text as an argument, like this:

```bash
# to execute the prompt in non-interactive mode
claude -p "center the div in index.html"
```

However, you need to give permissions to claude code to let it do things like edit files or use the Figma MCP server. 
Here are the flags to add:

```bash
# allows the MCP server called "figma-mcp-1" to be used
--allowedTools mcp__figma-mcp-1 

# allows claude to edit or create files
--permission-mode acceptEdits

# so, for example, this lets claude read a figma design and create the code for it:
claude -p "create this figma design in HTML/CSS: <url here>" --allowedTools mcp__figma-mcp-1 --permission-mode acceptEdits
```

> JSON Output:
> You can also get claude code to respond in JSON format. This is especially useful if you want to use claude code in a script.
> Read more about this here: https://docs.anthropic.com/en/docs/claude-code/sdk#json-output

### Script: Generate Designs for a List of URLs

I started with this, and it worked quite well actually. In my first claude code test, I noticed that there were a lot of inaccuracies and I had to iterate several times to start getting the design fully created. This made me wonder if I was feeding too much design data into claude code at once, so instead I tried getting the figma URLs for each top-level layer of the figma design. The results were significantly better.

So, I end up with a bunch of `index.html` and `styles.css` files for each top-level layer of the design.

The next part to consider is how to put them all together. There are two possibilities:

- I manually put them together. Not too much work to be honest, so totally doable.
- I try to get Cursor or Claude Code to piece them together. Not sure how this would work, but I like the idea of automating this process as much as possible.

Let's try asking claude code, in interactive mode, to look at the figma design and then try to put all the individual HTML files together into one big HTML file.

I gave claude code this prompt:

```
There is a directory called "generated", where within there are several directories that each
have the HTML/CSS for a layer in a figma design. I now want to put them all together into a single HTML file (and
single CSS file). Can you do this for me? Take a look at this figma design, and try to put all of these HTML files
together for me.

https://www.figma.com/design/AxOg233wgtoO27RFd8QIOv/M55_%E3%82%B3%E3%83%BC%E3%83%86%E3%82%99%E3%82%A3%E3%83%B3%E3%
82%AF%E3%82%99%E3%83%86%E3%82%B9%E3%83%88%E7%94%A8?node-id=1-1061&t=cXDUkWIRVTwE31dg-4
```

... and, it actually did the job more or less perfectly! It put together 6 different HTML pages - all representing different top-level layers from the main design - together, using the main design itself as reference.

Here's the cost of this "putting together" step, as output from the CLI:

```
Total cost:            $0.4831
Total duration (API):  6m 6.0s
```

## Using a `.fig` File to Generate HTML

I tried giving Claude Code a `.fig` file (file representing figma design data) to see if it can read the data and generate HTML from it.

Unfortunately, it seems like the `.fig` files may be encrypted or containing binary data, which Claude can't read:

```
⏺ Read(M55_コーディングテスト用.fig)
  ⎿  Error: This tool cannot read binary files. The file appears to be a
     binary .fig file. Please use appropriate tools for binary file
     analysis.
```

From here, Claude tried various ways to read URLs, File Keys, etc from the data in the file. It seemed to fail on each attempt:

```
⏺ The .fig file is a binary file that contains a Figma URL or file key. Let
   me try to extract the file key using bash to search for URL patterns.

⏺ Bash(strings "/Users/benwebb/dev/hexabase/hxai/temp/fig-file/M55_コーディ
      ングテスト用.fig" | grep -E "(figma\.com|file/)" | head -5)
  ⎿  (No content)

⏺ Let me try a different approach to find the file key:

⏺ Bash(strings "/Users/benwebb/dev/hexabase/hxai/temp/fig-file/M55_コーディ
      ングテスト用.fig" | grep -E "[a-zA-Z0-9]{22,}" | head -10)
  ⎿ nBBBBBBBBBBBBBBBBBBBBBBCC
    `gGGggGGggGGggGGggGGggGGggGGggGGggGGggGGggGGggGGggGGggGGggGGggGGggGGggG
    GggGGggGGggGGggGGggGGggGGggGGggGGggGGggGGggGGggGGggGGggGGggGG
    … +6 lines (ctrl+r to expand)


⏺ Let me try to find the actual Figma file key by looking for specific patterns:

⏺ Bash(hexdump -C
      "/Users/benwebb/dev/hexabase/hxai/temp/fig-file/M55_コーディングテスト用.fig" |
      grep -i "figma\|file\|key" | head -5)
  ⎿  (No content)

⏺ Let me try a different approach and look for the file key pattern more broadly:

⏺ Bash(strings
      "/Users/benwebb/dev/hexabase/hxai/temp/fig-file/M55_コーディングテスト用.fig" |
      grep -E "[A-Za-z0-9]{20,50}" | head -20)
  ⎿ nBBBBBBBBBBBBBBBBBBBBBBCC
    `gGGggGGggGGggGGggGGggGGggGGggGGggGGggGGggGGggGGggGGggGGggGGggGGggGGggGGggGGggGGggGG
    ggGGggGGggGGggGGggGGggGGggGGggGGggGGggGGggGGggGG
    … +6 lines (ctrl+r to expand)


⏺ It seems the .fig file doesn't contain a readable file key. Let me try a different
  approach and look for any URL patterns:

⏺ Bash(strings
      "/Users/benwebb/dev/hexabase/hxai/temp/fig-file/M55_コーディングテスト用.fig" |
      grep -i "http\|www\|\.com" | head -10)
  ⎿  yuwwwwwwwwwwwwwww
     pwwwwwwwwwwww
     qwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwww
     … +6 lines (ctrl+r to expand)
```

In the end, Claude admitted defeat and had this conclusion:

```
⏺ Unfortunately, the .fig file format doesn't contain a readable file key. The .fig
  files exported from Figma are usually just shortcuts that contain the URL to the Figma
   file, but this one appears to be corrupted or in a different format.

  To generate HTML/CSS from a Figma design, I would need either:
  1. The actual Figma file URL (like https://www.figma.com/file/[FILE_KEY]/...)
  2. A valid Figma file key
```

So, it seems that these files don't actually contain all the data needed to generate Figma designs anyway; Claude seems to think that they only could contain a URL or file key which it would then need to access directly from Figma's servers anyway. I think using a figma file directly is probably not a good option.

## Dev Notes

Here are some things I've learned about claude code as I've been working with the CLI and building my own wrapper CLI:

### MCP Server Scopes

When you register an MCP server, it might not be globally accessible (i.e. usable by claude code regardless of where you are in your directories). This depends on scope.

In the instructions above, we are setting a **local** scope. This means it's only accessible in that same directory.

If you want an MCP server that's globally accessible, then you need to use the **user** scope. More on that here:

https://docs.anthropic.com/en/docs/claude-code/mcp#user-scope