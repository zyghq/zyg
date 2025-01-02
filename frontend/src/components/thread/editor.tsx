import { Button } from "@/components/ui/button";
import { ButtonProps } from "@/components/ui/button";
import { sendThreadMailMessage } from "@/db/api.ts";
import { ThreadMessageResponse } from "@/db/schema.ts";
import { useMutation } from "@tanstack/react-query";
import CharacterCount from "@tiptap/extension-character-count";
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
import React from "react";

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

export function RichTextEditor({
  refetch,
  subject,
  threadId,
  token,
  workspaceId,
}: {
  refetch: () => void;
  subject: string;
  threadId: string;
  token: string;
  workspaceId: string;
}) {
  async function submit(html: string) {
    await mutation.mutateAsync({ htmlBody: html });
  }

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
            class: "pl-4 border-l-2 border-muted-foreground/40 italic my-1",
          },
        },
        bulletList: {
          HTMLAttributes: {
            class: "list-disc ml-4",
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
            class: "list-decimal ml-4",
          },
          itemTypeName: "listItem",
          keepMarks: true,
        },
        paragraph: {
          HTMLAttributes: {
            class: "leading-relaxed text-base text-gray-700",
          },
        },
      }),
      Placeholder.configure({
        emptyEditorClass: "is-editor-empty",
        placeholder: "Start typing...",
      }),
      Link.configure({
        HTMLAttributes: {
          class:
            "text-primary underline decoration-primary hover:text-primary/80 transition-colors",
        },
        openOnClick: false,
      }),
      CharacterCount.configure(),
    ],
  });

  const mutation = useMutation({
    mutationFn: async (values: { htmlBody: string }) => {
      const { htmlBody } = values;
      const { data, error } = await sendThreadMailMessage(
        token,
        workspaceId,
        threadId,
        { htmlBody },
      );
      if (error) {
        throw new Error(error.message);
      }
      if (!data) {
        throw new Error("no data returned");
      }
      return data as ThreadMessageResponse;
    },
    onError: (error) => {
      console.error(error);
    },
    onSuccess: (data) => {
      console.log("onSuccess");
      console.log(data);
      refetch();
      if (editor) editor.commands.clearContent();
    },
  });

  if (!editor) {
    return null;
  }

  const chars = editor.storage.characterCount.characters();

  return (
    <div className="flex h-full flex-col rounded-md border bg-white px-4 py-2 shadow-md dark:bg-background">
      <div className="flex items-center gap-2">
        <div className="flex-1 text-sm font-medium text-muted-foreground">
          {`Re: ${subject || ""}`}
        </div>
        <div className="flex items-center gap-2"></div>
      </div>
      <div className="flex-grow overflow-auto">
        <EditorContent className="h-full pr-2" editor={editor} />
      </div>
      <div className="flex items-center justify-between">
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
            disabled={chars === 0}
            onClick={() => submit(editor.getHTML())}
          >
            Send
          </Button>
        </div>
      </div>
      <div className="mt-2 flex h-1 items-center justify-end">
        {mutation.isError && (
          <div className="text-xs font-semibold text-red-600">
            Something went wrong
          </div>
        )}
      </div>
    </div>
  );
}
