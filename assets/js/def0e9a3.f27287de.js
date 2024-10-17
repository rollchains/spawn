"use strict";(self.webpackChunkdocs=self.webpackChunkdocs||[]).push([[104],{2776:(e,n,t)=>{t.r(n),t.d(n,{assets:()=>l,contentTitle:()=>a,default:()=>u,frontMatter:()=>r,metadata:()=>o,toc:()=>c});var s=t(5893),i=t(1151);const r={title:"Name Service",sidebar_label:"Testnet",sidebar_position:5,slug:"/build/name-service-testnet"},a="Running your Application",o={id:"build-your-application/testnet",title:"Name Service",description:"Congrats!! You built your first network already. You are ready to run a local testnet environment to verify it works.",source:"@site/versioned_docs/version-v0.50.x/02-build-your-application/05-testnet.md",sourceDirName:"02-build-your-application",slug:"/build/name-service-testnet",permalink:"/spawn/v0.50/build/name-service-testnet",draft:!1,unlisted:!1,tags:[],version:"v0.50.x",sidebarPosition:5,frontMatter:{title:"Name Service",sidebar_label:"Testnet",sidebar_position:5,slug:"/build/name-service-testnet"},sidebar:"defaultSidebar",previous:{title:"Configure Client",permalink:"/spawn/v0.50/build/name-service-client"},next:{title:"Bonus",permalink:"/spawn/v0.50/build/name-service-bonus"}},l={},c=[{value:"Launch The Network",id:"launch-the-network",level:3},{value:"Interact Set Name",id:"interact-set-name",level:3},{value:"Interaction Get Name",id:"interaction-get-name",level:2}];function d(e){const n={admonition:"admonition",code:"code",em:"em",h1:"h1",h2:"h2",h3:"h3",li:"li",p:"p",pre:"pre",ul:"ul",...(0,i.a)(),...e.components};return(0,s.jsxs)(s.Fragment,{children:[(0,s.jsx)(n.h1,{id:"running-your-application",children:"Running your Application"}),"\n",(0,s.jsxs)(n.admonition,{title:"Synopsis",type:"note",children:[(0,s.jsx)(n.p,{children:"Congrats!! You built your first network already. You are ready to run a local testnet environment to verify it works."}),(0,s.jsxs)(n.ul,{children:["\n",(0,s.jsx)(n.li,{children:"Building your application executable"}),"\n",(0,s.jsx)(n.li,{children:"Running a local testnet"}),"\n",(0,s.jsx)(n.li,{children:"Interacting with the network"}),"\n"]})]}),"\n",(0,s.jsx)(n.h3,{id:"launch-the-network",children:"Launch The Network"}),"\n",(0,s.jsxs)(n.p,{children:["Use the ",(0,s.jsx)(n.code,{children:"sh-testnet"})," command ",(0,s.jsx)(n.em,{children:"(short for shell testnet)"})," to quickly build your application, generate example wallet accounts, and start the local network on your machine."]}),"\n",(0,s.jsx)(n.pre,{children:(0,s.jsx)(n.code,{className:"language-bash",children:"# Run a quick shell testnet\nmake sh-testnet\n"})}),"\n",(0,s.jsx)(n.p,{children:"The chain will begin to create (mint) new blocks. You can see the logs of the network running in the terminal."}),"\n",(0,s.jsx)(n.h3,{id:"interact-set-name",children:"Interact Set Name"}),"\n",(0,s.jsxs)(n.p,{children:["Using the newly built binary ",(0,s.jsx)(n.em,{children:"(rolld from the --bin flag when the chain was created)"}),", you are going to execute the ",(0,s.jsx)(n.code,{children:"set"}),' transaction to your name. In this example, use "alice". This links account ',(0,s.jsx)(n.code,{children:"acc1"})," address to the desired name in the keeper."]}),"\n",(0,s.jsxs)(n.p,{children:["Then, resolve this name with the nameservice lookup. ",(0,s.jsx)(n.code,{children:"$(rolld keys show acc1 -a)"})," is a substitute for the acc1's address. You can also use just ",(0,s.jsx)(n.code,{children:"roll1hj5fveer5cjtn4wd6wstzugjfdxzl0xpg2te87"})," here."]}),"\n",(0,s.jsx)(n.pre,{children:(0,s.jsx)(n.code,{className:"language-bash",children:"rolld tx nameservice set alice --from=acc1 --yes\n\n# You can verify this transaction was successful\n# By querying it's unique ID.\nrolld q tx EC3FBF3248E24B5FEB6A5F7F35BBB4634E9C75587119E3FBCF5C1FED05E5A399\n"})}),"\n",(0,s.jsx)(n.h2,{id:"interaction-get-name",children:"Interaction Get Name"}),"\n",(0,s.jsxs)(n.p,{children:["Now you are going to get the name of a wallet. A nested command ",(0,s.jsx)(n.code,{children:"$(rolld keys show acc1 -a)"})," gets the unique address of the acc1 account added when you started the testnet."]}),"\n",(0,s.jsx)(n.pre,{children:(0,s.jsx)(n.code,{className:"language-bash",children:"rolld q nameservice resolve roll1efd63aw40lxf3n4mhf7dzhjkr453axur57cawh --output=json\n\nrolld q nameservice resolve $(rolld keys show acc1 -a) --output=json\n"})}),"\n",(0,s.jsx)(n.p,{children:"The expected result should be:"}),"\n",(0,s.jsx)(n.pre,{children:(0,s.jsx)(n.code,{className:"language-json",children:'{\n  "name": "alice"\n}\n'})}),"\n",(0,s.jsx)(n.admonition,{type:"note",children:(0,s.jsxs)(n.p,{children:["When you are ready to stop the testnet, you can use ",(0,s.jsx)(n.code,{children:"ctrl + c"})," or ",(0,s.jsx)(n.code,{children:"killall -9 rolld"}),"."]})}),"\n",(0,s.jsx)(n.p,{children:"Your network is now running and you have successfully set and resolved a name! \ud83c\udf89"})]})}function u(e={}){const{wrapper:n}={...(0,i.a)(),...e.components};return n?(0,s.jsx)(n,{...e,children:(0,s.jsx)(d,{...e})}):d(e)}},1151:(e,n,t)=>{t.d(n,{Z:()=>o,a:()=>a});var s=t(7294);const i={},r=s.createContext(i);function a(e){const n=s.useContext(r);return s.useMemo((function(){return"function"==typeof e?e(n):{...n,...e}}),[n,e])}function o(e){let n;return n=e.disableParentContext?"function"==typeof e.components?e.components(i):e.components||i:a(e.components),s.createElement(r.Provider,{value:n},e.children)}}}]);