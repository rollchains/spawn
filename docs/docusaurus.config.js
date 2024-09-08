// @ts-check
// Note: type annotations allow type checking and IDEs autocompletion

import { Highlight, themes } from "prism-react-renderer";

const lightCodeTheme = themes.github;
const darkCodeTheme = themes.dracula;

const organizationName = "rollchains";
const projectName = "spawn";

/** @type {import('@docusaurus/types').Config} */
const config = {
  title: "Spawn",
  tagline: "Documentation for Spawn",
  favicon: "img/white-cosmos-icon.svg",

  // TODO:
  // Set the production url of your site here
  // for local production tests, set to http://localhost:3000/
  // url: "https://ibc.cosmos.network",
  url: `https://${organizationName}.github.io`,
  // Set the /<baseUrl>/ pathname under which your site is served
  // For GitHub pages deployment, it is often '/<projectName>/'
  baseUrl: `/${projectName}/`,

  // GitHub pages deployment config.
  // If you aren't using GitHub pages, you don't need these.
  organizationName: organizationName, // Usually your GitHub org/user name.
  projectName: projectName, // Usually your repo name.
  deploymentBranch: "gh-pages",
  trailingSlash: true,

  onBrokenLinks: "log",
  onBrokenMarkdownLinks: "log",

  // Even if you don't use internalization, you can use this field to set useful
  // metadata like html lang. For example, if your site is Chinese, you may want
  // to replace "en" with "zh-Hans".
  i18n: {
    defaultLocale: "en",
    locales: ["en"],
  },

  presets: [
    [
      "classic",
      /** @type {import('@docusaurus/preset-classic').Options} */
      ({
        docs: {
          sidebarPath: require.resolve("./sidebars.js"),
          // Routed the docs to the root path
          routeBasePath: "/",
          // Exclude template markdown files from the docs
          exclude: ["**/*.template.md"],
          // Select the latest version
          lastVersion: "v0.50.x",
          // Assign banners to specific versions
          // editUrl: `https://github.com/${organizationName}/${projectName}/tree/main/`,
          versions: {
            current: {
              path: "main",
              banner: "unreleased",
            },
            "v0.50.x": {
              path: "v0.50",
              banner: "none",
            },
          },
        },
        theme: {
          customCss: require.resolve("./src/css/custom.css"),
        },
        // gtag: {
        //   trackingID: "G-HP8ZXWVLJG",
        //   anonymizeIP: true,
        // },
        sitemap: {
          changefreq: "weekly",
          priority: 0.5,
          filename: "sitemap.xml",
        },
      }),
    ],
  ],

  themeConfig:
    /** @type {import('@docusaurus/preset-classic').ThemeConfig} */
    ({
      // image: "img/ibc-go-docs-social-card.png",
      navbar: {
        logo: {
          alt: "Rollchains Logo",
          src: "img/Rollchains-Logo-RC-Black.svg",
          srcDark: "img/Rollchains-Logo-RC-White.svg",
          href: "/",
        },
        items: [
          // {
          //   type: "docSidebar",
          //   sidebarId: "defaultSidebar",
          //   position: "left",
          //   label: "Documentation",
          // },
          // {
          //   type: "doc",
          //   position: "left",
          //   docId: "README",
          //   docsPluginId: "adrs",
          //   label: "Architecture Decision Records",
          // },
          // {
          //   type: "doc",
          //   position: "left",
          //   docId: "intro",
          //   docsPluginId: "tutorials",
          //   label: "Tutorials",
          // },
          {
            type: "docsVersionDropdown",
            position: "right",
            dropdownActiveClassDisabled: true,
          },
          {
            href: "https://github.com/rollchains/spawn",
            html: `<svg width="24" height="24" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg" class="github-icon">
            <path fill-rule="evenodd" clip-rule="evenodd" d="M12 0.300049C5.4 0.300049 0 5.70005 0 12.3001C0 17.6001 3.4 22.1001 8.2 23.7001C8.8 23.8001 9 23.4001 9 23.1001C9 22.8001 9 22.1001 9 21.1001C5.7 21.8001 5 19.5001 5 19.5001C4.5 18.1001 3.7 17.7001 3.7 17.7001C2.5 17.0001 3.7 17.0001 3.7 17.0001C4.9 17.1001 5.5 18.2001 5.5 18.2001C6.6 20.0001 8.3 19.5001 9 19.2001C9.1 18.4001 9.4 17.9001 9.8 17.6001C7.1 17.3001 4.3 16.3001 4.3 11.7001C4.3 10.4001 4.8 9.30005 5.5 8.50005C5.5 8.10005 5 6.90005 5.7 5.30005C5.7 5.30005 6.7 5.00005 9 6.50005C10 6.20005 11 6.10005 12 6.10005C13 6.10005 14 6.20005 15 6.50005C17.3 4.90005 18.3 5.30005 18.3 5.30005C19 7.00005 18.5 8.20005 18.4 8.50005C19.2 9.30005 19.6 10.4001 19.6 11.7001C19.6 16.3001 16.8 17.3001 14.1 17.6001C14.5 18.0001 14.9 18.7001 14.9 19.8001C14.9 21.4001 14.9 22.7001 14.9 23.1001C14.9 23.4001 15.1 23.8001 15.7 23.7001C20.5 22.1001 23.9 17.6001 23.9 12.3001C24 5.70005 18.6 0.300049 12 0.300049Z" fill="currentColor"/>
            </svg>
            `,
            position: "right",
          },
        ],
      },
      footer: {
        links: [
          {
            // items: [
            //   {
            //     html: `<a href="https://cosmos.network"><img src="/img/cosmos-logo-bw.svg" alt="Cosmos Logo"></a>`,
            //   },
            // ],
          },
          // {
          //   title: "Documentation",
          //   items: [
          //     {
          //       label: "Hermes Relayer",
          //       href: "https://hermes.informal.systems/",
          //     },
          //     {
          //       label: "Cosmos Hub",
          //       href: "https://hub.cosmos.network",
          //     },
          //     {
          //       label: "CometBFT",
          //       href: "https://docs.cometbft.com",
          //     },
          //   ],
          // },
          {
            title: "Community",
            items: [
              {
                label: "Discord",
                href: "https://discord.com/invite/interchain",
              },
              {
                label: "Twitter",
                href: "https://x.com/rollchains",
              },
            ],
          },
          {
            title: "Other Tools",
            items: [
              {
                label: "interchaintest",
                href: "https://github.com/strangelove-ventures/interchaintest",
              },
              {
                label: "Go Relayer",
                href: "https://github.com/cosmos/relayer",
              },
              {
                label: "CosmWasm",
                href: "https://docs.cosmwasm.com/",
              },
            ],
          },
          {
            title: "More",
            items: [
              {
                label: "GitHub",
                href: "https://github.com/rollchains/spawn",
              },
              {
                label: "Rollchains Website",
                href: "https://www.rollchains.com",
              },
            ],
          },
        ],
        // logo: {
        //   alt: "Large IBC Logo",
        //   src: "img/black-large-ibc-logo.svg",
        //   srcDark: "img/white-large-ibc-logo.svg",
        //   width: 275,
        // },
      },
      prism: {
        theme: lightCodeTheme,
        darkTheme: darkCodeTheme,
        additionalLanguages: ["protobuf", "go", "go-module", "yaml", "toml", "diff", "bash"],
        magicComments: [
          // Remember to extend the default highlight class name as well!
          {
            className: 'theme-code-block-highlighted-line',
            line: 'highlight-next-line',
            block: { start: 'highlight-start', end: 'highlight-end' },
          },
          {
            className: 'code-block-minus-diff-line',
            line: 'minus-diff-line',
            block: { start: 'minus-diff-start', end: 'minus-diff-end' },
          },
          {
            className: 'code-block-plus-diff-line',
            line: 'plus-diff-line',
            block: { start: 'plus-diff-start', end: 'plus-diff-end' },
          },
        ],
      },
    }),
  themes: ["docusaurus-theme-github-codeblock"],
  plugins: [
    [
      'docusaurus-pushfeedback', {
        project: 'paq392agxs',
        buttonPosition: 'center-right',
        modalPosition: 'sidebar-right',
        buttonStyle: 'dark',
      }
    ],
    // [
    //   "@docusaurus/plugin-content-docs",
    //   {
    //     id: "tutorials",
    //     path: "tutorials",
    //     routeBasePath: "tutorials",
    //     sidebarPath: require.resolve("./sidebars.js"),
    //     exclude: ["**/*.template.md"],
    //   },
    // ],
    [
      "@docusaurus/plugin-client-redirects",
      {
        // makes the default page next in production
        redirects: [
          // {
          //   from: ["/master", "/next"],
          //   to: "/main/",
          // },
          {
            from: ["/", "/docs", "/spawn"],
            to: `/v0.50/`,
          }
        ],
      },
    ],
    [
      "@gracefullight/docusaurus-plugin-microsoft-clarity",
      { projectId: "idk9udvhuu" },
    ],
    [
      require.resolve("@easyops-cn/docusaurus-search-local"),
      {
        indexBlog: false,
        docsRouteBasePath: ["/"],
        highlightSearchTermsOnTargetPage: true,
      },
    ],
    async function myPlugin(context, options) {
      return {
        name: "docusaurus-tailwindcss",
        configurePostCss(postcssOptions) {
          postcssOptions.plugins.push(require("postcss-import"));
          postcssOptions.plugins.push(require("tailwindcss/nesting"));
          postcssOptions.plugins.push(require("tailwindcss"));
          postcssOptions.plugins.push(require("autoprefixer"));
          return postcssOptions;
        },
      };
    },
  ],
};

module.exports = config;
