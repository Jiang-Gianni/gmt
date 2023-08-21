# GMT

## Installation

```go
go install github.com/Jiang-Gianni/gmt
```



## How to use

This tool converts Markdown files into HTML.

If all the markdown files are in the same directory:
```bash
gmt -dir exampleMdDirectory -out exampleOutDirectory
```

In case of a single file:
```bash
gmt -dir exampleMdFile.md -out exampleOutDirectory
```

**-dir** and **-out** have both "." (current directory) as default.

If you want to add an external stylesheet link:
```bash
gmt -dir exampleMdDirectory -out exampleOutDirectory -css externalCss.css
```



## Tag Attributes

By writing a comment before any element inside the markdown file, the contents of the comment are injected into the HTML element during the parsing.

**Does not work inside code blocks**

```md
<!-- id="my-text" -->
**My Text**
```

is parsed as:

```html
<p  id="my-text" ><strong>My Text</strong></p>
```



## Tailwind

If Tailwind classes are used then a **style** tag is added on top of the converted HTML file containing the Tailwind CSS styling definitions.

```md
<!-- class="text-blue-700" -->
# My blue title
```

is parsed as:

```html
<style>.text-blue-700{--text-opacity:1;color:#2b6cb0;color:rgba(43,108,176,var(--text-opacity))}</style>
<h1 id="my-blue-title"  class="text-blue-700" >My blue title</h1>
```




## Templating

It can be used to generate html templates for various programming languages.

In order for the comment contents not to be parsed inside a **p** tag in the HTML, put the comment between two ticks **``**.

Using [**quicktemplate**](https://github.com/valyala/quicktemplate)'s syntax as an example:

```md
<!-- `{% func myGoFunction(name string) %}` -->
<!-- class="text-red-500" -->
**Hello <!-- `{%s name %}` -->**
<!-- `{% endfunc %}` -->
```

is parsed as:

```html
{% func myGoFunction(name string) %}
<p  class="text-red-500" ><strong>Hello {%s name %}</strong></p>
{% endfunc %}
```



## Pikchr

[**Pikchr**](https://pikchr.org/home/doc/trunk/homepage.md) support.

Pikchr's code between **\```pikchr** and **```** are converted into svg.

```pikchr
arrow right 200% "Markdown" "Source"
box rad 10px "Markdown" "Formatter" "(markdown.c)" fit
arrow right 200% "HTML+SVG" "Output"
arrow <-> down 70% from last box.s
box same "Pikchr" "Formatter" "(pikchr.c)" fit
```

is parsed as:

```html
<div class="pikchr-svg" style="max-width:423px"><svg xmlns='http://www.w3.org/2000/svg' viewBox="0 0 423.821 195.84">
<polygon points="146,37 134,41 134,33" style="fill:rgb(0,0,0)"/>
<path d="M2,37L140,37"  style="fill:none;stroke-width:2.16;stroke:rgb(0,0,0);" />
<text x="74" y="25" text-anchor="middle" fill="rgb(0,0,0)" dominant-baseline="central">Markdown</text>
<text x="74" y="49" text-anchor="middle" fill="rgb(0,0,0)" dominant-baseline="central">Source</text>
<path d="M161,72L258,72A15 15 0 0 0 273 57L273,17A15 15 0 0 0 258 2L161,2A15 15 0 0 0 146 17L146,57A15 15 0 0 0 161 72Z"  style="fill:none;stroke-width:2.16;stroke:rgb(0,0,0);" />
<text x="209" y="17" text-anchor="middle" fill="rgb(0,0,0)" dominant-baseline="central">Markdown</text>
<text x="209" y="37" text-anchor="middle" fill="rgb(0,0,0)" dominant-baseline="central">Formatter</text>
<text x="209" y="57" text-anchor="middle" fill="rgb(0,0,0)" dominant-baseline="central">(markdown.c)</text>
<polygon points="417,37 405,41 405,33" style="fill:rgb(0,0,0)"/>
<path d="M273,37L411,37"  style="fill:none;stroke-width:2.16;stroke:rgb(0,0,0);" />
<text x="345" y="25" text-anchor="middle" fill="rgb(0,0,0)" dominant-baseline="central">HTML+SVG</text>
<text x="345" y="49" text-anchor="middle" fill="rgb(0,0,0)" dominant-baseline="central">Output</text>
<polygon points="209,72 214,84 205,84" style="fill:rgb(0,0,0)"/>
<polygon points="209,123 205,111 214,111" style="fill:rgb(0,0,0)"/>
<path d="M209,78L209,117"  style="fill:none;stroke-width:2.16;stroke:rgb(0,0,0);" />
<path d="M176,193L243,193A15 15 0 0 0 258 178L258,138A15 15 0 0 0 243 123L176,123A15 15 0 0 0 161 138L161,178A15 15 0 0 0 176 193Z"  style="fill:none;stroke-width:2.16;stroke:rgb(0,0,0);" />
<text x="209" y="138" text-anchor="middle" fill="rgb(0,0,0)" dominant-baseline="central">Pikchr</text>
<text x="209" y="158" text-anchor="middle" fill="rgb(0,0,0)" dominant-baseline="central">Formatter</text>
<text x="209" y="178" text-anchor="middle" fill="rgb(0,0,0)" dominant-baseline="central">(pikchr.c)</text>
</svg>
</div>
```
![](./assets/pikchr.svg)
