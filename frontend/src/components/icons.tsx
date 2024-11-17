import { cn } from "@/lib/utils";
import { ChatBubbleIcon, DotFilledIcon } from "@radix-ui/react-icons";
import {
  CheckCircleIcon,
  ClockIcon,
  LocateIcon,
  MailIcon,
  PauseIcon,
  ReplyIcon,
} from "lucide-react";

export interface IconProps {
  className?: string;
  size?: number;
}

export const Icons = {
  logo: (props: IconProps) => (
    <svg
      height={40}
      viewBox="0 0 40 40"
      width={40}
      xmlns="http://www.w3.org/2000/svg"
      {...props}
    >
      <rect fill="black" height={40} rx={5} ry={5} width={40} />
      <text
        fill="white"
        fontFamily="sans-serif"
        fontSize={24}
        fontWeight="500"
        textAnchor="middle"
        x={22}
        y={28}
      >
        {"Z."}
      </text>
    </svg>
  ),
  nothing: (props: IconProps) => (
    <svg
      fill="none"
      viewBox="0 0 180 120"
      xmlns="http://www.w3.org/2000/svg"
      {...props}
    >
      <path
        d="M76.8492 20.7993C78.4007 19.2478 80.9214 19.2494 82.4794 20.8073C85.4239 23.7519 90.1943 23.7593 93.1343 20.8194C96.0743 17.8794 96.0668 13.109 93.1222 10.1645C91.5643 8.6065 91.5627 6.0858 93.1142 4.53432C94.6657 2.98285 97.1864 2.98439 98.7443 4.54235C104.805 10.6029 104.815 20.4146 98.7725 26.4575C92.7295 32.5005 82.9178 32.4899 76.8573 26.4294C75.2993 24.8715 75.2978 22.3508 76.8492 20.7993Z"
        stroke="#E3E4E8"
        strokeLinecap="round"
        strokeLinejoin="round"
        strokeWidth="1.25"
      />
      <path
        d="M143.43 70.9896C144.003 68.8519 146.675 68.1358 148.24 69.7008L153.532 74.9929C155.097 76.5579 154.381 79.2301 152.243 79.8029L145.014 81.74C142.876 82.3128 140.92 80.3566 141.493 78.2188L143.43 70.9896Z"
        stroke="#E3E4E8"
        strokeWidth="1.25"
      />
      <path
        d="M45.452 12.1193C47.5897 11.5465 49.5459 13.5027 48.9731 15.6404L48.286 18.2046C47.7132 20.3424 45.041 21.0584 43.4761 19.4935L41.5989 17.6163C40.034 16.0514 40.75 13.3792 42.8878 12.8064L45.452 12.1193Z"
        stroke="#E3E4E8"
        strokeWidth="1.25"
      />
      <path
        d="M72.8871 80.2408L72.8871 80.2455C72.8874 80.2843 72.8799 80.3227 72.865 80.3585C72.8501 80.3943 72.8282 80.4267 72.8005 80.4538C72.7728 80.4809 72.7399 80.5022 72.7038 80.5163C72.6698 80.5296 72.6335 80.5364 72.5969 80.5361C69.5409 80.3596 66.6008 79.304 64.1304 77.4961C61.6569 75.686 59.7597 73.1991 58.6677 70.3352C57.5756 67.4713 57.3353 64.3526 57.9756 61.3552C58.6159 58.3578 60.1095 55.6095 62.2764 53.4418C64.4433 51.2741 67.191 49.7794 70.1882 49.1379C73.1854 48.4965 76.3041 48.7357 79.1684 49.8267C82.0328 50.9177 84.5204 52.8139 86.3314 55.2867C88.1413 57.7581 89.1982 60.6999 89.3749 63.758C89.376 63.7952 89.3698 63.8324 89.3565 63.8672C89.3429 63.9031 89.322 63.9358 89.2952 63.9634C89.2684 63.9909 89.2363 64.0127 89.2008 64.0274C89.1653 64.042 89.1272 64.0493 89.0888 64.0487L89.0888 64.0486H89.0792H73.5121C73.1669 64.0486 72.8871 64.3285 72.8871 64.6736L72.8871 80.2408Z"
        stroke="#FFCC00"
        strokeLinejoin="round"
        strokeWidth="1.25"
      />
      <mask fill="white" id="path-5-inside-1_1614_85607">
        <path d="M102.068 76.5198C104.44 78.6363 107.367 80.0275 110.495 80.5257C113.624 81.024 116.821 80.608 119.702 79.3279C122.582 78.0478 125.022 75.9583 126.728 73.3112C128.435 70.6641 129.334 67.5725 129.318 64.4091C129.302 61.2457 128.371 58.1455 126.638 55.4823C124.905 52.8191 122.443 50.7065 119.55 49.3993C116.657 48.0922 113.455 47.6461 110.331 48.115C107.207 48.5839 104.295 49.9478 101.944 52.0421C101.856 52.1253 101.786 52.2256 101.738 52.3371C101.689 52.4485 101.664 52.5687 101.663 52.6906C101.663 52.8125 101.687 52.9334 101.734 53.0463C101.781 53.1592 101.85 53.2617 101.937 53.3477L112.921 64.3319L102.044 75.2094C101.956 75.2944 101.887 75.3966 101.84 75.5097C101.794 75.6228 101.771 75.7445 101.773 75.8672C101.776 75.9899 101.803 76.1111 101.854 76.2233C101.905 76.3355 101.978 76.4364 102.068 76.5198Z" />
      </mask>
      <path
        d="M101.944 52.0421L101.113 51.1088C101.104 51.1168 101.095 51.1248 101.086 51.133L101.944 52.0421ZM101.937 53.3477L102.821 52.4638L102.814 52.4572L101.937 53.3477ZM112.921 64.3319L113.805 65.2158C114.293 64.7276 114.293 63.9362 113.805 63.448L112.921 64.3319ZM102.044 75.2094L102.914 76.1069C102.918 76.1024 102.923 76.0979 102.928 76.0933L102.044 75.2094ZM101.84 75.5097L102.997 75.9849L101.84 75.5097ZM101.854 76.2233L102.993 75.7084L101.854 76.2233ZM101.236 77.4524C103.786 79.7277 106.933 81.2241 110.299 81.7602L110.692 79.2913C107.801 78.8308 105.094 77.5448 102.901 75.5871L101.236 77.4524ZM110.299 81.7602C113.665 82.2962 117.107 81.849 120.209 80.4702L119.194 78.1856C116.536 79.367 113.583 79.7517 110.692 79.2913L110.299 81.7602ZM120.209 80.4702C123.312 79.0913 125.941 76.8402 127.779 73.9884L125.678 72.634C124.104 75.0763 121.852 77.0043 119.194 78.1856L120.209 80.4702ZM127.779 73.9884C129.617 71.1367 130.585 67.8074 130.568 64.4028L128.068 64.4154C128.083 67.3376 127.252 70.1916 125.678 72.634L127.779 73.9884ZM130.568 64.4028C130.551 60.9983 129.549 57.6638 127.686 54.8005L125.59 56.1641C127.193 58.6272 128.053 61.493 128.068 64.4154L130.568 64.4028ZM127.686 54.8005C125.822 51.9372 123.176 49.6658 120.065 48.2602L119.035 50.5385C121.711 51.7472 123.987 53.7009 125.59 56.1641L127.686 54.8005ZM120.065 48.2602C116.953 46.8545 113.509 46.3741 110.146 46.8789L110.517 49.3512C113.401 48.9182 116.36 49.3298 119.035 50.5385L120.065 48.2602ZM110.146 46.8789C106.783 47.3837 103.645 48.8523 101.113 51.1088L102.776 52.9753C104.944 51.0432 107.632 49.7842 110.517 49.3512L110.146 46.8789ZM101.086 51.133C100.875 51.3328 100.706 51.5733 100.591 51.84L102.884 52.8341C102.865 52.878 102.838 52.9179 102.802 52.9511L101.086 51.133ZM100.591 51.84C100.475 52.1067 100.415 52.3937 100.413 52.6835L102.913 52.6976C102.913 52.7438 102.903 52.7903 102.884 52.8341L100.591 51.84ZM100.413 52.6835C100.412 52.9733 100.468 53.2601 100.58 53.527L102.887 52.5657C102.905 52.6068 102.914 52.6516 102.913 52.6976L100.413 52.6835ZM100.58 53.527C100.691 53.7938 100.854 54.0355 101.06 54.2381L102.814 52.4572C102.845 52.4879 102.87 52.5247 102.887 52.5657L100.58 53.527ZM101.053 54.2315L112.037 65.2158L113.805 63.448L102.821 52.4638L101.053 54.2315ZM112.037 63.448L101.16 74.3256L102.928 76.0933L113.805 65.2158L112.037 63.448ZM101.174 74.312C100.962 74.5168 100.796 74.7628 100.684 75.0346L102.997 75.9849C102.978 76.0305 102.95 76.0721 102.914 76.1069L101.174 74.312ZM100.684 75.0346C100.572 75.3064 100.518 75.5978 100.524 75.8908L103.023 75.8436C103.024 75.8912 103.015 75.9393 102.997 75.9849L100.684 75.0346ZM100.524 75.8908C100.529 76.1837 100.594 76.472 100.715 76.7382L102.993 75.7084C103.012 75.7502 103.022 75.7961 103.023 75.8436L100.524 75.8908ZM100.715 76.7382C100.835 77.0044 101.008 77.2434 101.223 77.4407L102.914 75.5988C102.947 75.6294 102.974 75.6666 102.993 75.7084L100.715 76.7382Z"
        fill="#FFCC00"
        mask="url(#path-5-inside-1_1614_85607)"
      />
      <rect
        height="9.31628"
        rx="4.65814"
        stroke="#E3E4E8"
        strokeLinecap="round"
        strokeLinejoin="round"
        strokeWidth="1.25"
        transform="rotate(-15 18.1395 39.0874)"
        width="27.3706"
        x="18.1395"
        y="39.0874"
      />
      <rect
        height="9.31628"
        rx="4.65814"
        stroke="#E3E4E8"
        strokeLinecap="round"
        strokeLinejoin="round"
        strokeWidth="1.25"
        transform="rotate(-15 0.765466 66.8687)"
        width="41.0151"
        x="0.765466"
        y="66.8687"
      />
      <rect
        height="7.80681"
        rx="3.90341"
        stroke="#FF7D26"
        strokeWidth="1.25"
        width="7.80681"
        x="56.1045"
        y="28.5566"
      />
      <rect
        height="7.80681"
        rx="3.90341"
        stroke="#FF7D26"
        strokeWidth="1.25"
        width="7.80681"
        x="30.4434"
        y="76.8623"
      />
      <rect
        height="7.80681"
        rx="3.90341"
        stroke="#FF7D26"
        strokeWidth="1.25"
        width="7.80681"
        x="144.401"
        y="94.2207"
      />
      <rect
        height="7.80681"
        rx="3.90341"
        stroke="#E3E4E8"
        strokeLinecap="round"
        strokeLinejoin="round"
        strokeWidth="1.25"
        width="7.80681"
        x="171.567"
        y="30.0684"
      />
      <rect
        height="9.31628"
        rx="4.65814"
        stroke="#E3E4E8"
        strokeLinecap="round"
        strokeLinejoin="round"
        strokeWidth="1.25"
        transform="rotate(-15 113.63 22.9654)"
        width="34.9772"
        x="113.63"
        y="22.9654"
      />
      <rect
        height="9.31628"
        rx="4.65814"
        stroke="#E3E4E8"
        strokeLinecap="round"
        strokeLinejoin="round"
        strokeWidth="1.25"
        transform="rotate(-15 155.982 57.2378)"
        width="21.392"
        x="155.982"
        y="57.2378"
      />
      <rect
        height="9.31628"
        rx="4.65814"
        stroke="#E3E4E8"
        strokeLinecap="round"
        strokeLinejoin="round"
        strokeWidth="1.25"
        transform="rotate(-15 134.042 40.6939)"
        width="27.4299"
        x="134.042"
        y="40.6939"
      />
      <rect
        height="9.31628"
        rx="4.65814"
        stroke="#E3E4E8"
        strokeLinecap="round"
        strokeLinejoin="round"
        strokeWidth="1.25"
        transform="rotate(-15 22.284 97.1499)"
        width="21.392"
        x="22.284"
        y="97.1499"
      />
      <path
        d="M75.9228 112.849C74.6238 115.099 71.3762 115.099 70.0772 112.849L57.9528 91.8486C56.6538 89.5986 58.2776 86.7861 60.8756 86.7861L85.1244 86.7861C87.7224 86.7861 89.3462 89.5986 88.0472 91.8486L75.9228 112.849Z"
        stroke="#2BD95A"
        strokeWidth="1.25"
      />
      <path
        d="M117.923 112.849C116.624 115.099 113.376 115.099 112.077 112.849L99.9528 91.8486C98.6538 89.5986 100.278 86.7861 102.876 86.7861L127.124 86.7861C129.722 86.7861 131.346 89.5986 130.047 91.8486L117.923 112.849Z"
        stroke="#45B4FF"
        strokeWidth="1.25"
      />
    </svg>
  ),
  oops: (props: IconProps) => (
    <svg viewBox="0 0 70 70" xmlns="http://www.w3.org/2000/svg" {...props}>
      <path
        className="color000000 svgShape"
        d="M3.6894531,62.6342773C5.7661133,65.9941406,9.3632812,68,13.3129883,68h43.3740234
	c3.949707,0,7.546875-2.0058594,9.6235352-5.3657227c2.0761719-3.359375,2.2612305-7.4741211,0.4951172-11.0068359
	L45.1186523,8.2539062C43.1899414,4.3964844,39.3129883,2,35,2s-8.1899414,2.3964844-10.1186523,6.2539062L3.1943359,51.6274414
	C1.4282227,55.1601562,1.6132812,59.2749023,3.6894531,62.6342773z M4.9833984,52.5219727L26.6704102,9.1484375
	C28.2822266,5.9248047,31.3959961,4,35,4s6.7177734,1.9248047,8.3295898,5.1484375l21.6870117,43.3735352
	c1.4541016,2.9077148,1.3017578,6.2954102-0.4077148,9.0610352C62.8999023,64.3486328,59.9384766,66,56.6870117,66H13.3129883
	c-3.2514648,0-6.2128906-1.6513672-7.921875-4.4169922C3.6816406,58.8173828,3.5292969,55.4296875,4.9833984,52.5219727z"
        fill="#d85b53"
      ></path>
      <path
        className="color000000 svgShape"
        d="M34.9995117 47.3867188c2.6943359 0 4.8867188-2.1918945 4.8867188-4.8862305V23.1547852c0-2.6943359-2.1923828-4.8862305-4.8867188-4.8862305s-4.8862305 2.1918945-4.8862305 4.8862305v19.3457031C30.1132812 45.1948242 32.3051758 47.3867188 34.9995117 47.3867188zM32.1132812 23.1547852c0-1.5913086 1.2949219-2.8862305 2.8862305-2.8862305 1.5917969 0 2.8867188 1.2949219 2.8867188 2.8862305v19.3457031c0 1.5913086-1.2949219 2.8862305-2.8867188 2.8862305-1.5913086 0-2.8862305-1.2949219-2.8862305-2.8862305V23.1547852zM35 59.4702148c2.7568359 0 5-2.2431641 5-5s-2.2431641-5-5-5-5 2.2431641-5 5S32.2431641 59.4702148 35 59.4702148zM35 51.4702148c1.6542969 0 3 1.3457031 3 3s-1.3457031 3-3 3-3-1.3457031-3-3S33.3457031 51.4702148 35 51.4702148z"
        fill="#d85b53"
      ></path>
    </svg>
  ),
  spinner: ({ size = 24, ...props }) => (
    <svg
      height={size}
      width={size}
      xmlns="http://www.w3.org/2000/svg"
      {...props}
      className={cn("animate-spin", props.className)}
      fill="none"
      stroke="currentColor"
      strokeLinecap="round"
      strokeLinejoin="round"
      strokeWidth="2"
      viewBox="0 0 24 24"
    >
      <path d="M21 12a9 9 0 1 1-6.219-8.56" />
    </svg>
  ),
};

export const PriorityIcons = {
  high: (props: IconProps) => (
    <svg
      fill="none"
      height="18"
      viewBox="0 0 18 18"
      width="18"
      xmlns="http://www.w3.org/2000/svg"
      {...props}
    >
      <rect fill="#f5b458" height="2" rx="0.5" width="12" x="3" y="4"></rect>
      <rect fill="#f5b458" height="2" rx="0.5" width="12" x="3" y="8"></rect>
      <rect fill="#f5b458" height="2" rx="0.5" width="12" x="3" y="12"></rect>
    </svg>
  ),
  low: (props: IconProps) => (
    <svg
      fill="none"
      height="18"
      viewBox="0 0 18 18"
      width="18"
      xmlns="http://www.w3.org/2000/svg"
      {...props}
    >
      <rect fill="#e4e4ea" height="2" rx="0.5" width="12" x="3" y="4"></rect>
      <rect fill="#e4e4ea" height="2" rx="0.5" width="12" x="3" y="8"></rect>
      <rect fill="#9898a9" height="2" rx="0.5" width="12" x="3" y="12"></rect>
    </svg>
  ),
  normal: (props: IconProps) => (
    <svg
      fill="none"
      height="18"
      viewBox="0 0 18 18"
      width="18"
      xmlns="http://www.w3.org/2000/svg"
      {...props}
    >
      <rect fill="#e4e4ea" height="2" rx="0.5" width="12" x="3" y="4"></rect>
      <rect fill="#9898a9" height="2" rx="0.5" width="12" x="3" y="8"></rect>
      <rect fill="#9898a9" height="2" rx="0.5" width="12" x="3" y="12"></rect>
    </svg>
  ),
  urgent: (props: IconProps) => (
    <svg
      fill="none"
      height="18"
      viewBox="0 0 18 18"
      width="18"
      xmlns="http://www.w3.org/2000/svg"
      {...props}
    >
      <path
        clipRule="evenodd"
        d="M2.674 3.778C2 4.787 2 6.19 2 9c0 2.809 0 4.213.674 5.222.292.437.667.812 1.104 1.104C4.787 16 6.19 16 9 16c2.809 0 4.213 0 5.222-.674a4.003 4.003 0 0 0 1.104-1.104C16 13.213 16 11.81 16 9c0-2.809 0-4.213-.674-5.222a4.002 4.002 0 0 0-1.104-1.104C13.213 2 11.81 2 9 2c-2.809 0-4.213 0-5.222.674a4 4 0 0 0-1.104 1.104ZM9 10c.283 0 .52-.096.713-.287A.968.968 0 0 0 10 9V6a.968.968 0 0 0-.287-.713A.968.968 0 0 0 9 5a.968.968 0 0 0-.713.287A.968.968 0 0 0 8 6v3c0 .283.096.52.287.713.192.191.43.287.713.287Zm0 3c.283 0 .52-.096.713-.287A.968.968 0 0 0 10 12a.968.968 0 0 0-.287-.713A.967.967 0 0 0 9 11a.967.967 0 0 0-.713.287A.968.968 0 0 0 8 12c0 .283.096.52.287.713.192.191.43.287.713.287Z"
        fill="#e33d3d"
        fillRule="evenodd"
      ></path>
    </svg>
  ),
};

// TODO: add other icons for `spam` or `ignored`
export function stageIcon(stage: string, props: IconProps) {
  switch (stage) {
    case "hold":
      return <PauseIcon {...props} />;
    case "needs_first_response":
      return <LocateIcon {...props} />;
    case "needs_next_response":
      return <ReplyIcon {...props} />;
    case "resolved":
      return <CheckCircleIcon {...props} />;
    case "waiting_on_customer":
      return <ClockIcon {...props} />;
    default:
      return <></>;
  }
}

export function channelIcon(channel: string, props: IconProps) {
  switch (channel) {
    case "chat":
      return <ChatBubbleIcon {...props} />;
    case "email":
      return <MailIcon {...props} />;
    default:
      return <></>;
  }
}

export function eventSeverityIcon(severity: string, props: IconProps) {
  switch (severity) {
    case "critical":
      return (
        <DotFilledIcon
          {...props}
          className={cn(
            props.className,
            "animate-pulse text-red-600 dark:text-red-400",
          )}
        />
      );
    case "error":
      return (
        <DotFilledIcon
          {...props}
          className={cn(
            props.className,
            "animate-pulse text-rose-600 dark:text-rose-400",
          )}
        />
      );
    case "info":
      return (
        <DotFilledIcon
          {...props}
          className={cn(props.className, "text-blue-600 dark:text-blue-400")}
        />
      );
    case "muted":
      return (
        <DotFilledIcon
          {...props}
          className={cn(props.className, "text-gray-500 dark:text-gray-400")}
        />
      );
    case "success":
      return (
        <DotFilledIcon
          {...props}
          className={cn(props.className, "text-green-600 dark:text-green-400")}
        />
      );
    case "warning":
      return (
        <DotFilledIcon
          {...props}
          className={cn(
            props.className,
            "text-yellow-600 dark:text-yellow-400",
          )}
        />
      );
    default:
      return <DotFilledIcon {...props} />;
  }
}
