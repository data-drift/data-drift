import type { Meta, StoryObj } from "@storybook/react";
import { CommitListItem } from "./CommitListItem";

const meta = {
  title: "Drift/Commits/CommitListItem",
  component: CommitListItem,
} satisfies Meta<typeof CommitListItem>;

export default meta;

type Story = StoryObj<typeof meta>;

export const SimpleCase: Story = {};
