import type { Meta, StoryObj } from "@storybook/react";
import { CommitListItem } from "./CommitListItem";

const meta = {
  title: "Drift/Commits/CommitListItem",
  component: CommitListItem,
} satisfies Meta<typeof CommitListItem>;

export default meta;

type Story = StoryObj<typeof meta>;

export const DriftCase: Story = {
  args: {
    type: "Drift",
    isParentData: true,
    date: new Date("2021-08-01T00:00:00Z"),
    commitUrl: "www.google.com",
    name: "New drift detected data/act_metrics_finance/bop_eop_mrr_monthly_by_country.csv (#76)\n\nDrift: data/act_metrics_finance/bop_eop_mrr_monthly_by_country.csv",
  },
};

export const NewDataCase: Story = {
  args: {
    type: "New Data",
    isParentData: false,
    date: new Date("2021-08-01T00:00:00Z"),
    commitUrl: "www.google.com",
    name: "New drift detected data/act_metrics_finance/bop_eop_mrr_monthly_by_country.csv (#76)\n\nDrift: data/act_metrics_finance/bop_eop_mrr_monthly_by_country.csv",
  },
};
