/**
 * Creating a sidebar enables you to:
 - create an ordered group of docs
 - render a sidebar for each doc of that group
 - provide next/previous navigation

 The sidebars can be generated from the filesystem, or explicitly defined here.

 Create as many sidebars as you want.
 */

// @ts-check

/** @type {import('@docusaurus/plugin-content-docs').SidebarsConfig} */
const sidebars = {
  // By default, Docusaurus generates a sidebar from the docs folder structure
  defaultSidebar: [
    { type: "autogenerated", dirName: "." },
    // {
    //   type: "category",
    //   label: "Resources",
    //   collapsed: true,
    //   items: [
    //     {
    //       type: "link",
    //       label: "IBC",
    //       href: "https://ibc.cosmos.network/v8",
    //     },
    //     {
    //       type: "link",
    //       label: "Cosmos-SDK",
    //       href: "https://docs.cosmos.network",
    //     },
    //     {
    //       type: "link",
    //       label: "Developer Portal",
    //       href: "https://tutorials.cosmos.network",
    //     },
    //     {
    //       type: "link",
    //       label: "Awesome Cosmos",
    //       href: "https://github.com/cosmos/awesome-cosmos",
    //     },
    //   ],
    // },
  ],

  // But you can create a sidebar manually
  /*
  tutorialSidebar: [
    'intro',
    'hello',
    {
      type: 'category',
      label: 'Tutorial',
      items: ['tutorial-basics/create-a-document'],
    },
  ],
   */
};

module.exports = sidebars;
