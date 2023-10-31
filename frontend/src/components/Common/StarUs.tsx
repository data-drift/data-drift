import GitHubButton from "react-github-btn";

const StarUs = () => (
  <GitHubButton
    href="https://github.com/data-drift/data-drift"
    data-color-scheme="no-preference: dark; light: light; dark: dark;"
    data-icon="octicon-star"
    aria-label="Star data-drift/data-drift on GitHub"
    data-size="large"
    data-show-count={true}
  >
    Star
  </GitHubButton>
);

export default StarUs;
