import { Badge } from "@/components/ui/badge";
import { buttonVariants } from "@/components/ui/button";
import { Button } from "@/components/ui/button";
import { Separator } from "@/components/ui/separator";
import {
  Tooltip,
  TooltipContent,
  TooltipProvider,
  TooltipTrigger,
} from "@/components/ui/tooltip";
import { cn } from "@/lib/utils";
import {
  CheckCircledIcon,
  CopyIcon,
  DotFilledIcon,
} from "@radix-ui/react-icons";
import { useCopyToClipboard } from "@uidotdev/usehooks";
import Markdown from "react-markdown";

type TextSize = "L" | "M" | "S" | "XS";

type TextColor = "ERROR" | "MUTED" | "NORMAL" | "SUCCESS" | "WARNING";

interface ComponentTextProps {
  text: string;
  textColor?: TextColor;
  textSize?: TextSize;
}

// Renders Markdown text component.
export function ComponentText({
  text,
  textColor = "NORMAL",
  textSize = "S",
}: ComponentTextProps) {

  const sizeMap: Record<TextSize, string> = {
    L: "text-lg",
    M: "text-base",
    S: "text-sm",
    XS: "text-xs",
  };

  const colorMap: Record<TextColor, string> = {
    ERROR: "text-red-600 dark:text-red-400",
    MUTED: "text-muted-foreground",
    NORMAL: "text-gray-900 dark:text-gray-100",
    SUCCESS: "text-green-600 dark:text-green-400",
    WARNING: "text-yellow-600 dark:text-yellow-400",
  };

  return (
    <Markdown
      components={{
        code: ({ children }) => (
          <code className="block whitespace-pre-wrap break-words">
            {children}
          </code>
        ),
        h1: ({ children }) => <h1 className="text-base">{children}</h1>,
        h2: ({ children }) => <h2 className="text-base">{children}</h2>,
        h3: ({ children }) => <h3 className="text-base">{children}</h3>,
        h4: ({ children }) => <h4 className="text-base">{children}</h4>,
        p: ({ children }) => (
          <p
            className={cn(sizeMap[textSize] || "text-sm", colorMap[textColor])}
          >
            {children}
          </p>
        ),
        pre: ({ children }) => (
          <pre className="max-w-full overflow-auto whitespace-pre-wrap break-words">
            {children}
          </pre>
        ),
      }}
    >
      {text}
    </Markdown>
  );
}

type SpacerSize = "L" | "M" | "S" | "XS";

interface ComponentSpacerProps {
  spacerSize: SpacerSize;
}

// Renders a spacer component.
export function ComponentSpacer({ spacerSize }: ComponentSpacerProps) {
  const sizeMap: Record<SpacerSize, string> = {
    L: "h-8",
    M: "h-4",
    S: "h-2",
    XS: "h-1",
  };
  return <div className={sizeMap[spacerSize] || "h-2"} />;
}

interface ComponentLinkButtonProps {
  linkButtonLabel: string;
  linkButtonUrl: string;
}

// Renders a link button component.
export function ComponentLinkButton({
  linkButtonLabel,
  linkButtonUrl,
}: ComponentLinkButtonProps) {
  return (
    <a
      className={buttonVariants({ variant: "outline" })}
      href={linkButtonUrl}
      rel="noopener noreferrer"
      target="_blank"
    >
      {linkButtonLabel}
    </a>
  );
}

// Renders plain text component.
export function ComponentPlainText({
  text,
  textSize = "S",
}: ComponentTextProps) {
  const sizeMap: Record<TextSize, string> = {
    L: "text-lg",
    M: "text-base",
    S: "text-sm",
    XS: "text-xs",
  };

  return <p className={sizeMap[textSize] || "text-sm"}>{text}</p>;
}

type DividerSpacingSize = "L" | "M" | "S" | "XS";

interface ComponentDividerProps {
  dividerSize?: DividerSpacingSize;
}

// Renders a divider component.
export function ComponentDivider({ dividerSize = "S" }: ComponentDividerProps) {
  const sizeMap: Record<DividerSpacingSize, string> = {
    L: "my-8",
    M: "my-4",
    S: "my-2",
    XS: "my-1",
  };
  return <Separator className={sizeMap[dividerSize] || "my-2"} />;
}

interface ComponentCopyButtonProps {
  copyButtonToolTipLabel: string;
  copyButtonValue: string;
}

// Renders a copy button component.
export function ComponentCopyButton({
  copyButtonToolTipLabel,
  copyButtonValue,
}: ComponentCopyButtonProps) {
  const [copiedText, copyToClipboard] = useCopyToClipboard();
  const hasCopiedText = Boolean(copiedText);

  return (
    <TooltipProvider>
      <Tooltip>
        <TooltipTrigger asChild>
          <Button
            className="text-muted-foreground"
            onClick={() => copyToClipboard(copyButtonValue || "")}
            size="icon"
            type="button"
            variant="ghost"
          >
            {hasCopiedText ? (
              <CheckCircledIcon className="h-4 w-4 text-green-600 dark:text-green-400" />
            ) : (
              <CopyIcon className="h-4 w-4" />
            )}
          </Button>
        </TooltipTrigger>
        <TooltipContent>
          <p>{copyButtonToolTipLabel}</p>
        </TooltipContent>
      </Tooltip>
    </TooltipProvider>
  );
}

type BadgeColor = "BLUE" | "GRAY" | "GREEN" | "RED" | "YELLOW";

interface ComponentBadgeProps {
  badgeColor?: BadgeColor;
  badgeLabel: string;
}

export function ComponentBadge({
  badgeColor = "BLUE",
  badgeLabel,
}: ComponentBadgeProps) {
  const dotColorMap: Record<BadgeColor, string> = {
    BLUE: "text-blue-600 dark:text-blue-300",
    GRAY: "text-gray-600 dark:text-gray-300",
    GREEN: "text-green-600 dark:text-green-300",
    RED: "text-red-600 dark:text-red-300",
    YELLOW: "text-yellow-600 dark:text-yellow-300",
  };

  const badgeColorMap: Record<BadgeColor, string> = {
    BLUE: "bg-blue-50 dark:bg-blue-900/80 border border-blue-200 dark:border-blue-700",
    GRAY: "bg-gray-50 dark:bg-gray-900/80 border border-gray-200 dark:border-gray-700",
    GREEN:
      "bg-green-50 dark:bg-green-900/80 border border-green-200 dark:border-green-700",
    RED: "bg-red-50 dark:bg-red-900/80 border border-red-200 dark:border-red-700",
    YELLOW:
      "bg-yellow-50 dark:bg-yellow-900/80 border border-yellow-200 dark:border-yellow-700",
  };

  const badgeColorCls = badgeColorMap[badgeColor] || badgeColorMap.GRAY;
  const dotColorCls = dotColorMap[badgeColor] || dotColorMap.GRAY;

  return (
    <Badge className={badgeColorCls} variant="outline">
      <DotFilledIcon className={cn(dotColorCls, "h-5 w-5")} />
      {badgeLabel}
    </Badge>
  );
}

type ComponentMap = {
  [key: string]: React.ComponentType<any>;
};

type ComponentProps = Record<string, any>;

interface RenderComponentsProps {
  components: ComponentProps[];
}

interface RowComponentProps {
  rowAsideContent: ComponentProps[];
  rowMainContent: ComponentProps[];
}

export function ComponentRow({
  rowAsideContent = [],
  rowMainContent = [],
}: RowComponentProps) {
  return (
    <div className="flex justify-between">
      {RenderComponents({ components: rowMainContent })}
      {RenderComponents({ components: rowAsideContent })}
    </div>
  );
}

const componentMap: ComponentMap = {
  componentBadge: ComponentBadge,
  componentCopyButton: ComponentCopyButton,
  componentDivider: ComponentDivider,
  componentLinkButton: ComponentLinkButton,
  componentPlainText: ComponentPlainText,
  componentRow: ComponentRow,
  componentSpacer: ComponentSpacer,
  componentText: ComponentText,
};

export function RenderComponents({ components = [] }: RenderComponentsProps) {
  const renderComponent = (component: ComponentProps, index: number) => {
    const [componentType, props]: [string, ComponentProps] =
      Object.entries(component)[0];

    const Component = componentMap[componentType];

    if (!Component) {
      console.warn(`Unknown component type: ${componentType}`);
      return null;
    }
    return <Component key={index} {...props} />;
  };

  return (
    <div>
      {components.map((component, index) => renderComponent(component, index))}
    </div>
  );
}
