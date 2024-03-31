export const Icons = {
  logo: (props) => (
    <svg
      xmlns="http://www.w3.org/2000/svg"
      width={40}
      height={40}
      viewBox="0 0 40 40"
      {...props}
    >
      <rect width={40} height={40} rx={5} ry={5} fill="black" />
      <text
        x={22}
        y={28}
        fontSize={24}
        fontFamily="sans-serif"
        fontWeight="500"
        textAnchor="middle"
        fill="white"
      >
        {"Z."}
      </text>
    </svg>
  ),
};
