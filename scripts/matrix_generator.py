# flake8:noqa
# Generate bash command matrixes for manual verification testing (due to the number of repos and testing required)
# rm -rf ./matrix*/
# cd scripts && python3 matrix_generator.py

import os
import random
import string

ConsensusFeatures = ["proof-of-authority", "proof-of-stake", "interchain-security"]

SupportedFeatures = [
    "tokenfactory",
    "ibc-packetforward",
    "ibc-ratelimit",
    "cosmwasm",
    "wasm-light-client",
    "optimistic-execution",
]

POS = "proof-of-stake"
POA = "proof-of-authority"
ICS = "interchain-security"  # cosmoshub


def main():
    current_dir = os.path.dirname(os.path.realpath(__file__))
    print(f"Generating matrix.sh in {current_dir}")

    cmds: list[str] = [
        #
        #
        "make install\n# rm -rf ./matrix*/\n\n# GH; Base - none disabled",
        CmdCreator("poabase", POA, [], "poa", "poad", "upoa", "strangelove")
        .set_push_to_gh()
        .build(),
        CmdCreator("posbase", POS, [], "cosmos", "appd", "ustake", "rollchains")
        .set_push_to_gh()
        .build(),
        CmdCreator("icsbase", ICS, [], "ics", "icsd", "uics", "my-ics")
        .set_push_to_gh()
        .build(),
        #
        #
        "\n# GH; Minimal generation (i.e. no features used)",
        CmdCreator(
            "spawndefaultfeatures",
            POA,
            "ibc-ratelimit,cosmwasm,wasm-light-client".split(","),
            "cosmos",
            "appd",
            "utoken",
            "myghorg",
        )
        .set_push_to_gh()
        .build(),
        CmdCreator(
            "posminimal",
            POS,
            SupportedFeatures,
            "minimal",
            "minid",
            "umini",
            "gh_org",
        )
        .set_push_to_gh()
        .build(),
        CmdCreator(
            "icsminimal",
            ICS,
            SupportedFeatures,
            "minimal",
            "minid",
            "umini",
            "gh_org",
        )
        .set_push_to_gh()
        .build(),
        #
        #
        "\n# Mixes - Ensure proper and local testing",
        CmdCreator(
            "randmixone",
            POS,
            "wasmlc".split(","),
            random_string(5, True),
            f"{random_string(6, True)}",
            f"u{random_string(5, True)}",
            random_string(10, True),
        )
        .set_with_local_unit_test()
        .build(),
        CmdCreator(
            "randmixtwo",
            POA,
            "wasmlc,tokenfactory".split(","),
            random_string(5, True),
            f"{random_string(6, True)}",
            f"u{random_string(5, True)}",
            random_string(10, True),
        )
        .set_with_local_unit_test()
        .build(),
        CmdCreator(
            "randmixthree",
            ICS,
            "packetforward,ibc-ratelimit".split(","),
            random_string(5, True),
            f"{random_string(6, True)}",
            f"u{random_string(5, True)}",
            random_string(10, True),
        )
        .set_with_local_unit_test()
        .build(),
        CmdCreator(
            "randmixfour",
            POA,
            "wasmlc,packetforward,tokenfactory".split(","),
            random_string(5, True),
            f"{random_string(6, True)}",
            f"u{random_string(5, True)}",
            random_string(10, True),
        )
        .set_with_local_unit_test()
        .build(),
        CmdCreator(
            "randmixfive",
            POS,
            "ibc-ratelimit".split(","),
            random_string(5, True),
            f"{random_string(6, True)}",
            f"u{random_string(5, True)}",
            random_string(10, True),
        )
        .set_with_local_unit_test()
        .build(),
        CmdCreator(
            "mywlcnocw",
            POS,
            "tokenfactory,ibc-packetforward,cosmwasm".split(","),
            random_string(5, True),
            f"{random_string(6, True)}",
            f"u{random_string(5, True)}",
            random_string(10, True),
        )
        .set_with_local_unit_test()
        .build(),
        #
        #
        "\n# Custom modules",
        CmdCreator(
            "poswithmodules",
            POS,
            [],
            random_string(5, True),
            f"{random_string(6, True)}",
            f"u{random_string(5, True)}",
            random_string(10, True),
        )
        .set_custom_modules(["aaaaa", "bbbbb"])
        .set_with_local_unit_test()
        .build(),
        CmdCreator(
            "poawithmoduleslong",
            POA,
            [],
            random_string(5, True),
            f"{random_string(6, True)}",
            f"u{random_string(5, True)}",
            random_string(10, True),
        )
        .set_custom_modules(
            ["mycns", "mysuperlongnameheremysuperlongnameheremysuperlongnamehere"]
        )
        .set_with_local_unit_test()
        .build(),
        CmdCreator(
            "icswithsomanymodules",
            ICS,
            [],
            random_string(5, True),
            f"{random_string(6, True)}",
            f"u{random_string(5, True)}",
            random_string(10, True),
        )
        .set_custom_modules([f"{random_string(6, True)}" for _ in range(5)] + [f"{random_string(6, True)} --ibc-module" for _ in range(5)] + [f"{random_string(6, True)} --ibc-middleware" for _ in range(10)])
        .set_with_local_unit_test()
        .build(),
        CmdCreator(
            "poswithwithibctest",
            POS,
            [],
            random_string(5, True),
            f"{random_string(6, True)}",
            f"u{random_string(5, True)}",
            random_string(10, True),
        )
        .set_custom_modules(["nsibc --ibc-module"])
        .set_with_local_unit_test()
        .build(),
    ]

    with open(f"{current_dir}/matrix.sh", "w") as f:
        f.write("\n".join(cmds))
        print("matrix.sh created")


class CmdCreator:
    name: str
    consensus: str
    disabled_features: list[str]
    bech32: str
    bin: str
    denom: str
    org: str
    open_in_code: bool
    push_to_gh: bool

    def __init__(
        self,
        name: str,
        consensus: str,
        features: list,
        bech32: str,
        bin: str,
        denom: str,
        org: str,
    ):
        self.name = f"matrix{name}"
        self.consensus = consensus
        self.disabled_features = features
        self.bech32 = bech32
        self.bin = bin
        self.denom = denom
        self.org = org

        self.go_mod_tidy = True
        self.open_in_code = False
        self.push_to_gh = False
        self.with_local_unit_test = False
        self.with_docker_build = False
        self.custom_modules: list[str] = []

    def set_open_in_code(self, open_in_code: bool) -> "CmdCreator":
        self.open_in_code = open_in_code
        return self

    def set_with_local_unit_test(self) -> "CmdCreator":
        self.with_local_unit_test = True
        return self

    def set_go_mod_tidy(self, go_mod_tidy: bool) -> "CmdCreator":
        self.go_mod_tidy = go_mod_tidy
        return self

    def set_with_docker_build(self) -> "CmdCreator":
        self.with_docker_build = True
        return self

    def set_custom_modules(self, modules: list[str]) -> "CmdCreator":
        self.custom_modules = modules
        return self

    def with_custom_module(self, module_name: str) -> "CmdCreator":
        self.custom_modules.append(module_name)
        return self

    def set_push_to_gh(self) -> "CmdCreator":
        self.push_to_gh = True
        return self

    def build(self) -> str:
        if len(self.name) == 0:
            self.name = random_string(8, True)
        if len(self.bech32) == 0:
            self.bech32 = random_string(4, True)
        if len(self.bin) == 0:
            self.bin = random_string(5, True) + "d"
        if len(self.denom) == 0:
            self.denom = "u" + random_string(3, True)
        if len(self.org) == 0:
            self.org = random_string(10, True)

        text = f"spawn new {self.name} --consensus={self.consensus} --bech32={self.bech32} --bin={self.bin} --denom={self.denom} --org={self.org}"

        if len(self.disabled_features) == 0:
            text += " --bypass-prompt"
            text += " --disable=explorer"
        else:
            text += " --disable=" + ",".join(self.disabled_features + ["explorer"])

        # if len(extraBashArgs) > 0:
        #     text += f" && {extraBashArgs}"

        if len(self.custom_modules) > 0:
            for module in self.custom_modules:
                text += f" && cd {self.name} && spawn module new {module} && cd .."
            text += f" && cd {self.name} && make proto-gen && cd .."

        if self.go_mod_tidy:
            text += f" && cd {self.name} && make mod-tidy && cd .."

        if self.open_in_code:
            text += f" && code {self.name}/"

        if self.with_local_unit_test:
            text += f" && cd {self.name} && go test ./... && cd .."

        if self.with_docker_build:
            text += f" && cd {self.name} && make local-image && cd .."

        if self.push_to_gh:
            # this is run after a cd .. so we set the source as that nested dir
            text += f" && gh repo create {self.name} --source={self.name}/ --remote=upstream --push --private"
            text += f"\ngh repo delete {self.name} --yes"

        return text


def random_string(length: int, alphaOnly: bool):
    text = string.ascii_letters
    if not alphaOnly:
        text += string.digits
    return "".join(random.choices(text, k=length))


if __name__ == "__main__":
    main()
