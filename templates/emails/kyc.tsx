import {
  Body,
  Button,
  Container,
  Head,
  Html,
  Img,
  Preview,
  Section,
  Text,
} from "@react-email/components";
import * as React from "react";

export const KycEmail = () => (
  <Html>
    <Head />
    <Preview>{`{{ .PreviewText }}`}</Preview>
    <Body style={main}>
      <Container style={container}>
        <Section>
          <Img
            style={logoImage}
            src="https://assets.zyg.ai/zyg.png"
            width="32"
            height="32"
            alt="Zyg"
          />
        </Section>
        <Text style={title}>
          <strong>You</strong> started a conversation.
        </Text>
        <Section style={section}>
          <Text style={text}>{`{{ .PreviewText }}`}</Text>
          <Button href={`{{ .MagicLink }}`} style={button}>
            Verify your email
          </Button>
        </Section>

        <Text style={footer}>
          ❤️ Zyg ・ Open source, made with love around the world ❤️
        </Text>
      </Container>
    </Body>
  </Html>
);

export default KycEmail;

const main = {
  backgroundColor: "#ffffff",
  color: "#24292e",
  fontFamily:
    '-apple-system,BlinkMacSystemFont,"Segoe UI",Helvetica,Arial,sans-serif,"Apple Color Emoji","Segoe UI Emoji"',
};

const container = {
  maxWidth: "480px",
  margin: "0 auto",
  padding: "20px 0 48px",
};

const logoImage = {
  marginTop: "0px",
  marginBottom: "0px",
  marginLeft: "auto",
  marginRight: "auto",
};

// const heading = {
//   marginTop: "12px",
//   marginBottom: "12px",
//   color: "#6a737d",
//   fontSize: "16px",
// };

const title = {
  fontSize: "16px",
  lineHeight: 1.25,
};

const section = {
  padding: "24px",
  border: "solid 1px #dedede",
  borderRadius: "5px",
  textAlign: "center" as const,
};

const text = {
  margin: "0 0 10px 0",
  textAlign: "left" as const,
};

const button = {
  fontWeight: "600",
  fontSize: "14px",
  backgroundColor: "#000000",
  color: "#fff",
  lineHeight: 1.5,
  borderRadius: "0.5em",
  padding: "10px 10px",
};

const footer = {
  color: "#6a737d",
  fontSize: "12px",
  textAlign: "center" as const,
  marginTop: "40px",
};
