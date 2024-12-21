import { Button } from "@/components/ui/button";
import { ButtonProps } from "@/components/ui/button";
import Link from "@tiptap/extension-link";
import Placeholder from "@tiptap/extension-placeholder";
import { EditorContent, useEditor } from "@tiptap/react";
import StarterKit from "@tiptap/starter-kit";
import {
  Bold,
  Code,
  Italic,
  List,
  ListOrdered,
  Strikethrough,
} from "lucide-react";

interface MenuButtonProps extends ButtonProps {
  active: boolean;
  children: React.ReactNode;
  onClick: () => void;
}

const MenuButton = ({
  active,
  children,
  onClick,
  ...props
}: MenuButtonProps): JSX.Element => (
  <Button
    className="h-7 w-7 hover:bg-muted data-[active=true]:bg-muted"
    data-active={active}
    onClick={onClick}
    size="icon"
    variant="ghost"
    {...props}
  >
    {children}
  </Button>
);

export function RichTextEditor() {
  const editor = useEditor({
    editorProps: {
      attributes: {
        class:
          "prose prose-sm dark:prose-invert prose-p:my-2 prose-pre:my-2 max-w-none focus:outline-none min-h-full",
      },
    },
    extensions: [
      StarterKit.configure({
        blockquote: {
          HTMLAttributes: {
            class:
              "zyg-blockquote pl-4 border-l-2 border-muted-foreground/40 italic my-2",
          },
        },
        bulletList: {
          HTMLAttributes: {
            class: "zyg-bulletlist list-disc ml-4",
          },
          itemTypeName: "listItem",
          keepMarks: true,
        },
        dropcursor: false,
        gapcursor: false,
        heading: false,
        horizontalRule: false,
        orderedList: {
          HTMLAttributes: {
            class: "zyg-orderedlist list-decimal ml-4",
          },
          itemTypeName: "listItem",
          keepMarks: true,
        },
        paragraph: {
          HTMLAttributes: {
            class: "zyg-paragraph leading-7",
          },
        },
      }),
      Placeholder.configure({
        emptyEditorClass: "is-editor-empty",
        placeholder: "Press R to reply...",
      }),
      Link.configure({
        HTMLAttributes: {
          class:
            "text-primary underline decoration-primary hover:text-primary/80 transition-colors",
        },
        openOnClick: false,
      }),
    ],
  });

  if (!editor) {
    return null;
  }

  return (
    <div className="flex h-full flex-col rounded-md border bg-white p-4 shadow-md dark:bg-background">
      <div className="flex items-center gap-2">
        <input
          className="flex-1 border-none bg-transparent text-sm font-medium text-muted-foreground outline-none"
          defaultValue="Re: Check the attachment for Bug"
          placeholder="Subject"
          type="text"
        />
        <div className="flex items-center gap-2">
          <Button size="sm" variant="ghost">
            Cc
          </Button>
          <Button size="sm" variant="ghost">
            Bcc
          </Button>
        </div>
      </div>
      <div className="flex-grow overflow-auto">
        <EditorContent className="h-full pr-2" editor={editor} />
      </div>
      <div className="flex items-center justify-between pt-2">
        <div className="flex items-center gap-1">
          <MenuButton
            active={editor.isActive("bold")}
            onClick={() => editor.chain().focus().toggleBold().run()}
          >
            <Bold className="h-4 w-4" />
            <span className="sr-only">Bold</span>
          </MenuButton>

          <MenuButton
            active={editor.isActive("italic")}
            onClick={() => editor.chain().focus().toggleItalic().run()}
          >
            <Italic className="h-4 w-4" />
            <span className="sr-only">Italic</span>
          </MenuButton>

          <MenuButton
            active={editor.isActive("strike")}
            onClick={() => editor.chain().focus().toggleStrike().run()}
          >
            <Strikethrough className="h-4 w-4" />
            <span className="sr-only">Strikethrough</span>
          </MenuButton>

          <MenuButton
            active={editor.isActive("bulletList")}
            onClick={() => editor.chain().focus().toggleBulletList().run()}
          >
            <List className="h-4 w-4" />
            <span className="sr-only">Bullet List</span>
          </MenuButton>

          <MenuButton
            active={editor.isActive("orderedList")}
            onClick={() => editor.chain().focus().toggleOrderedList().run()}
          >
            <ListOrdered className="h-4 w-4" />
            <span className="sr-only">Numbered List</span>
          </MenuButton>

          <MenuButton
            active={editor.isActive("code")}
            onClick={() => editor.chain().focus().toggleCode().run()}
          >
            <Code className="h-4 w-4" />
            <span className="sr-only">Code</span>
          </MenuButton>
        </div>

        <div className="flex items-center gap-2">
          <Button
            className="gap-2"
            onClick={() => {
              console.log("Submit content:", editor.getHTML());
            }}
            size="sm"
          >
            Reply & Investigate
            <kbd className="pointer-events-none inline-flex select-none items-center gap-1 rounded border bg-muted px-1.5 font-mono text-[10px] text-muted-foreground">
              <span>⌘</span>⏎
            </kbd>
          </Button>
        </div>
      </div>
    </div>
  );
}
