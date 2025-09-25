import { cookies } from "next/headers";
import { AppSidebar } from "@/components/app-sidebar";
import { ModeToggle } from "@/components/theme-toggle";
import {
  Breadcrumb,
  BreadcrumbItem,
  BreadcrumbLink,
  BreadcrumbList,
  BreadcrumbPage,
  BreadcrumbSeparator,
} from "@/components/ui/breadcrumb";
import { Separator } from "@/components/ui/separator";
import {
  SidebarInset,
  SidebarProvider,
  SidebarTrigger,
} from "@/components/ui/sidebar";
import MarkdownRenderer from "@/components/Markdown/MarkdownRenderer";

export default async function Page() {
  const cookieStore = await cookies();
  const defaultOpen = cookieStore.get("sidebar_state")?.value === "true";

  const markdown = `
  # Markdown Test

  ## Emphasis
  Some **bold text**, some *italic text*, and some ***bold italic text***.
  You can also ~~strikethrough~~ text.

  ## Lists

  ### Unordered List
  - Item 1
  - Item 2
    - Subitem 2a
    - Subitem 2b
  - Item 3

  ### Ordered List
  1. First
  2. Second
  3. Third

  ## Links
  - [OpenAI](https://www.openai.com)
  - Inline link: [Google](https://www.google.com)

  ## Images
  ![Sample Image](https://via.placeholder.com/150)

  ## Code

  ### Inline Code
  Here is some \`inline code\` within a sentence.

  ### Code Block
  \`\`\`javascript
  console.log("Hello, world!");
  function add(a, b) {
    return a + b;
  }
  \`\`\`

  ### Syntax Highlighting
  \`\`\`python
  def greet(name):
      return f"Hello, {name}!"
  \`\`\`

  ## Blockquotes
  > This is a blockquote.
  > It can span multiple lines.
  >> Nested blockquote example.

  ## Tables
  | Name       | Age | Occupation   |
  |------------|-----|--------------|
  | Alice      | 25  | Engineer     |
  | Bob        | 30  | Designer     |
  | Charlie    | 22  | Student      |

  ## Horizontal Rule
  ---
  Another section after the horizontal rule.

  ## Task List
  - [x] Completed task
  - [ ] Incomplete task
  - [ ] Another task

  `;


  return (
    <SidebarProvider defaultOpen={defaultOpen}>
      <AppSidebar />
      <SidebarInset className="border-1 border-border">
        <header className="flex h-16 shrink-0 items-center gap-2">
          <div className="flex items-center gap-2 px-4">
            <SidebarTrigger className="-ml-1" />
            <Separator
              orientation="vertical"
              className="mr-2 data-[orientation=vertical]:h-4"
            />
            <Breadcrumb>
              <BreadcrumbList>
                <BreadcrumbItem className="hidden md:block">
                  <BreadcrumbLink href="#">
                    Objektorientert Programmering
                  </BreadcrumbLink>
                </BreadcrumbItem>
                <BreadcrumbSeparator className="hidden md:block" />
                <BreadcrumbItem>
                  <BreadcrumbPage>Overview</BreadcrumbPage>
                </BreadcrumbItem>
              </BreadcrumbList>
            </Breadcrumb>
          </div>
          <div className="ml-auto px-4">
            <ModeToggle />
          </div>
        </header>
        <MarkdownRenderer content={markdown} />
      </SidebarInset>
    </SidebarProvider>
  );
}
