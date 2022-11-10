// SPDX-License-Identifier: MIT
//https://testnet.bscscan.com/address/0x1ba491c5078a109685877f5a368dde0303a83fe5#code
pragma solidity 0.8.17;
pragma experimental ABIEncoderV2;

contract Example {
    struct EIP712Domain {
        string  name;
        string  version;
        uint256 chainId;
        address verifyingContract;
    }

    struct Person {
        string name;
        address wallet;
    }

    struct Mail {
        Person from;
        Person to;
        string contents;
        uint256 amount;
        uint256 expiration;
    }

    bytes32 constant EIP712DOMAIN_TYPEHASH = keccak256(
        "EIP712Domain(string name,string version,uint256 chainId,address verifyingContract)"
    );

    bytes32 constant PERSON_TYPEHASH = keccak256(
        "Person(string name,address wallet)"
    );

    bytes32 constant MAIL_TYPEHASH = keccak256(
        "Mail(Person from,Person to,string contents,uint256 amount,uint256 expiration)Person(string name,address wallet)"
    );

    bytes32 public DOMAIN_SEPARATOR;

    constructor () {
        DOMAIN_SEPARATOR = hash(EIP712Domain({
        name: "Ether Mail",
        version: '1',
        chainId: 1,
        // verifyingContract: this
        verifyingContract: 0xCcCCccccCCCCcCCCCCCcCcCccCcCCCcCcccccccC
        }));

    }

    function hash(EIP712Domain memory eip712Domain) public pure returns (bytes32) {
        return keccak256(abi.encode(
                EIP712DOMAIN_TYPEHASH,
                keccak256(bytes(eip712Domain.name)),
                keccak256(bytes(eip712Domain.version)),
                eip712Domain.chainId,
                eip712Domain.verifyingContract
            ));
    }

    function hash(Person memory person) public pure returns (bytes32) {
        return keccak256(abi.encode(
                PERSON_TYPEHASH,
                keccak256(bytes(person.name)),
                person.wallet
            ));
    }

    function hash(Mail memory mail) public pure returns (bytes32) {
        return keccak256(abi.encode(
                MAIL_TYPEHASH,
                hash(mail.from),
                hash(mail.to),
                keccak256(bytes(mail.contents)),
                mail.amount,
                mail.expiration
            ));
    }

    function verify(Mail memory mail, uint8 v, bytes32 r, bytes32 s) public view returns (bool) {
        // Note: we need to use `encodePacked` here instead of `encode`.
        bytes32 digest = keccak256(abi.encodePacked(
                "\x19\x01",
                DOMAIN_SEPARATOR,
                hash(mail)
            ));
        return ecrecover(digest, v, r, s) == mail.from.wallet;
    }

    function verify(Mail memory mail, bytes memory signature) public view returns (bool) {
        // Note: we need to use `encodePacked` here instead of `encode`.
        bytes32 digest = keccak256(abi.encodePacked(
                "\x19\x01",
                DOMAIN_SEPARATOR,
                hash(mail)
            ));
        (uint8 v, bytes32 r, bytes32 s) = parseSignature(signature);
        return ecrecover(digest, v, r, s) == mail.from.wallet;
    }


    function parseSignature(bytes memory signature)
    public
    pure
    returns (
        uint8 v,
        bytes32 r,
        bytes32 s
    )
    {

        // The signature format is a compact form of:
        //   {bytes32 r}{bytes32 s}{uint8 v}
        // Compact means, uint8 is not padded to 32 bytes.
        assembly {
        // solium-disable-line security/no-inline-assembly
            r := mload(add(signature, 32))
            s := mload(add(signature, 64))
        // Here we are loading the last 32 bytes, including 31 bytes
        // of 's'. There is no 'mload8' to do this.
        //
        // 'byte' is not working due to the Solidity parser, so lets
        // use the second best option, 'and'
            v := and(mload(add(signature, 65)), 0xff)
        }

        if (v < 27) {
            v += 27;
        }
        require(v == 27 || v == 28, "invalid v of signature(r, s, v)");
    }

    function test() public view returns (bool) {
        // Example signed message
        Mail memory mail = Mail({
        from: Person({
        name: "Cow",
        wallet: 0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266
        }),
        to: Person({
        name: "Bob",
        wallet: 0xbBbBBBBbbBBBbbbBbbBbbbbBBbBbbbbBbBbbBBbB
        }),
        contents: "Hello, Bob!",
        amount: 10000000000,
        expiration: 1667659989
        });

        uint8 v = 27;
        bytes32 r = 0xfe1408ef873223d473a9130ecafad79c297b79286e38d90bb107f049831bbb61;
        bytes32 s = 0x11058bb70448776840216e6ceb27fe9f56a7fb9682fc63c0ab708b7bba296baa;

        assert(DOMAIN_SEPARATOR == 0xf2cee375fa42b42143804025fc449deafd50cc031ca257e0b194a650a912090f);
        assert(hash(mail) == 0x83fb919a6723739a9187fa6b145d321d7a747703fb65ff02e5adca18a3537c2a);
        assert(verify(mail, v, r, s));
        return true;
    }

    function getDigest() public view returns (bytes32){
        Mail memory mail = Mail({
        from: Person({
        name: "Cow",
        wallet: 0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266
        }),
        to: Person({
        name: "Bob",
        wallet: 0xbBbBBBBbbBBBbbbBbbBbbbbBBbBbbbbBbBbbBBbB
        }),
        contents: "Hello, Bob!",
        amount: 10000000000,
        expiration: 1667659989
        });

        bytes32 digest = keccak256(abi.encodePacked(
                "\x19\x01",
                DOMAIN_SEPARATOR,
                hash(mail)

            ));

        return digest;
    }
}