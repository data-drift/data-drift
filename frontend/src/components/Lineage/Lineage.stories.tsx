import { Meta, StoryObj } from "@storybook/react";
import Lineage from "./Lineage";

const meta = {
  title: "Lineage/Lineage",
  component: Lineage,
} satisfies Meta<typeof Lineage>;

export default meta;

type Story = StoryObj<typeof meta>;

export const DefaultCase: Story = {};
