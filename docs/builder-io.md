# Using Builder.io to Convert Figma to HTML

This guide will tell you how to generate HTML/CSS/JS code for a given figma design.

## The Manual Way

> You can find some detailed notes about this process in this Github issue, which also includes some screenshots:
>
> https://github.com/b-eee/sakai-ai/issues/1#issuecomment-3044541281

Once you've created a Builder.io account and created a space, you will be able to start generating code for a figma design.

To do this, go to the design of your choice in Figma and do the following:

1. Open the Builder.io Plugin

> This requires you to be in "Dev Mode", as does any other plugin in Figma

2. On the Export tab, choose to export your code

- you can choose either "Smart Export" or "Classic Export". I suppose there are some differences between the two, but both seem to work for me.
- However, in order to specifically generate HTML/CSS (rather than React), you may need to choose "Classic Export" for now. We will go this route for these instructions.

3. Once the export process has finished, you will have some options for how to proceed. Choose to "Open in Builder".

- The export process may take a few minutes, depending on how big and complex the design is.
- It may ask you to choose a preview version you like. Choose the one that looks closest to the original design.

4. Once Builder has opened, you will have some options to generate code. Customize the output to be in HTML and CSS (instead of React or other SPA frameworks).

- The option for customizing code may be found in the top right corner of the screen.

5. Generate the code, and view the results!

## The CLI Way

To get the code into a project more efficiently, you can use the command line instead. You'll still need to do steps 1 and 2 of the manual instructions, but it should make things a little faster overall.

1. Do steps 1-2 of the Manual instructions above.

> TODO: should we still do Classic export? Or is Smart Export better?

2. Once the export process has finished, copy the `npx` command. In your terminal, at the appropriate directory, paste this command and answer the prompts.

- If you paste it in an empty directory it will ask you if you want to continue or not. Answer "Yes".
- It will ask you how to implement the design into your project. This is where you can tell the AI what to do.
  - I typically tell it here to use HTML/CSS/JS and no SPA frameworks, in order to get plain HTML without anything like React.
  - However, if you configure things right, this may not be necessary (see "Configuring Rules for AI" sectionn below).

3. Once it's done, the CLI should tell you a summary of the changes and the files will have been created for you.

### Configuring Rules for AI

There are a couple files you can create to configure how Builder.io generates the code for you.
However, there isn't much documentation about how these _actually_ are formatted.

1. `.builderrules`: Use this to add custom instructions that will be automatically injected into the LLM prompt during code generation.
I'd assume this can just be natural language, plain text?

2. `.builderignore`: Use this file to exclude specific files or patterns from being included in the code generation process. This is especially useful for template files that hsouldn't influence the generated code, legacy code that doesn't follow the current standards, or third-party code that might confuse the generation process, and test files or mock data.
I assume that this just works similarly to `.gitignore` files, where files that match the given pattern are just not noticed by the AI agent?

So, for our purposes, I believe we can define a `.builderrules` file to tell the LLM to generate code that is as close to WebRelease2 stardards as possible.

> TODO: come up with a prompt that will get the code output as close to WebRelease2's HTML format as possible.

### CLI flags and options

When I try the `--help` flag, I see this information about the CLI command we use here:

```
$ npx "@builder.io/dev-tools@latest" code --help

┌   help   1.6.148
│
●  → Code Generation
│
│   npx builder.io code --url URL
│
│  Generate or modify code based on Figma designs.
│
│  Arguments:
│    --url            URL to start completion from
│    --spaceId        Builder.io space ID to use
│    --prompt         Prompt text for non-interactive mode
│    --workspace      Specify a workspace configuration file for multi-root workspaces
│    --mode           Generation mode - either 'exact' for precise matches or 'creative' for more flexibility
│    --cwd            Working directory to run commands from
│
│  Configuration Files:
│    .builderignore   Add patterns to exclude files from being included in code generation
│    .builderrules    Add custom instructions that will be injected into the LLM prompt
│    .cursorrules     Cursor's settings are automatically supported for consistent behavior
```

The most interesting one here is `--prompt`, which seems to let you set the prompt for the AI to use when generating code, so you don't have to enter it.
Still, since we currently have to copy the command from Figma's website, it's not that easy to add in a `--prompt` flag. Maybe we could create an alias in the shell that adds in a default prompt flag with the value we want?

Unfortunately, it doesn't look like there's a flag for specifying a `.fig` file.

## Benefit of using with Cursor?

This article goes into depth about how to generate code with builder.io and then reap benefits from Cursor:

https://www.builder.io/blog/figma-to-cursor

I've read through this, and it seems like the only point to using Cursor here is for **after** generating the code with Builder.io.
Cursor, as with any project, is able to analyze the code and help with making changes, refactoring, etc. So, while it's certainly helpful, it doesn't seem to be important for the actual initial code generation that builder.io does.

If you do use cursor, you can additionally make use of a `.cursorrules` file to tell Cursor's agent what kind of coding patterns to use.