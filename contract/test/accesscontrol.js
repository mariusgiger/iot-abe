const AccessControl = artifacts.require("AccessControl");
const assert = require("assert");
const truffleAssert = require("truffle-assertions");
const ethCrypto = require("eth-crypto");

contract("Test AccessControl Contract", accounts => {
  let acc;
  let owner = accounts[0];
  let requester = accounts[1];
  let iotDevice = accounts[2];
  let requesterPrivKey =
    "0x428b9597ab91a53c63fe06b3998c3519d4881452362c03e5ac202022784cae76";

  let pubKey =
    "7479706520610a7120383738303731303739393636333331323532323433373738313938343735343034393831353830363838333139393431343230383231313032383635333339393236363437353633303838303232323935373037383632353137393432323636323232313432333135353835383736393538323331373435393237373731333336373331373438313332343932353132393939383232343739310a682031323031363031323236343839313134363037393338383832313336363734303533343230343830323935343430313235313331313832323931393631353133313034373230373238393335393730343533313130323834343830323138333930363533373738363737360a72203733303735303831383636353435313632313336313131393234353537313530343930313430353937363535393631370a65787032203135390a65787031203130370a7369676e3120310a7369676e3020310a0000000080a454954ea2240a277bf67a767582e9cf1849e93f84ba391947f7026b10977c1e2bde1aaf83a9df571499eeea64f251da9972b4f99611135ea8b3be6658ec5eaa78f9821dd51a86b3221937bac4e65b3507a6f893000b9543e77cdb3bb78a06de457dc63ff276c8659ef37280e0f1207e3cd85d584e217fe2a290bed55240984a0000008060bbbbbd481a87ea17f40df23c6cd68bc4834b625e9720d1f8838a3159e811460b48419a3497cc8fef402953cff1f0a1abd50e78ab0c879d947d175298a5951e142376ac2f99d7f7f3db2281977066024b82aa5c9cee9533413937ac7e453e42ce531d6b9dc910da6ab5565ccc5bb4c1b1194deca5c04247cffdf00792b4285100000080375daf417eb317f93eb07903948a71b567ec0be4f266aaebbea4e482973b0c5b17d018ba36f0416079a69be6c7439e5ca9083118c73bd2f36d52efcfd9f8cbc537bf8eb3c254c70aa87b89180a58698877f6e902821cae4d702d5d9bca47caf729e9002a6e007ee374e44e778f790e06530cb59bccf12902c1c4bab34f58e40c0000008050082c3178c2b8367b37f039cbfd21a9e90bf9bc25bf315209e8d6a6e083f09bac8980494a4f68c023734c54926f3e4e52c3cea32b00f83be1e1155e4c21c8ee6ea2a35fa831fcad0852bb83d4171f27809630af9a3d73bb23db85b013c3723a1c8cea639727d78cdbbc0ec0623e7dbb4d3ff5483dde71b273afb4c38d9da32f";

  let masterKey =
    "000000147501846e73d948dc40fdc5d5f72dd0b22c2744de000000801b5c319fd90bd41d67f7125c1b758d7aa68cb59c8f2daf895d3ef4bf995cc47723d6376c1e6af396e227d58ce53ead6f3f526e9bb0d5e371c66f59ca43736aed04836d6016cb77aec467cf966020e171e46417f6f57ccf2571d61f0080a2874d0e81503b021ed99f737bc1f4786340bfb7c1199e19ba8aa9f9ddb3a34b026923";

  let prvKey =
    "0000008052fee8772a29d3a905e166e01e253243cc1300fdf6c372338cb62eb0805e5168d5f02a650cb91ac096cbe21f39056a683a0d8f887095a33e1169c32eee18d554039ed6a7183ebe3fdf7acb49b1a4e673c4dea3f4f0563f0a970ec56133716f89b328eca5ee7ba037db428651d97268a0415af3e90fb5d87948525f5ad7357d0e0000000261646d696e00000000805a8428903bc8ed6d4ea173d7d51b39ca295cb7b02dd21480c18a820fcb6521a229e0b6afbb0c4285eee1a7db13d8af706588b78b5f5821c9e9e725fa99c32e5e3bb0795479b5e7f7145f00041af1bc4c166524d9150ed2a0afbf55cee08b9aaaf5c0eaade18b00c7bbb7ede48efe787a729d5d35130691e608f361517c62062d0000008018661f2cd156ed3fdd004004c2d5a32aa995a40af00dfe9352c4fb92865d6cbc086243998cc8eb7037fab63d72db06aacfb266f1301e6325cd12a1a2818e3a110dd422cae6a4ed224577a77b1cf5b05e5bffa31f5805546e12206a8649b9faf3a0db179bb9a3f6f6c701bbccd3f0d2d9a9da40ab5b02b35008a5f21a3c0ab9d669745f646570617274656d656e74000000008041e312527b72263ae8e03e6a1b7365afba3e4f9278bb2fdfbdad0216938bf8e2ce8f648547627f82f595944a42849e9dd17f278044d1844610ce09b42415870959475db874749596165fb1dfc967625e43f561abea199a376a2933ad4b2377e6f0ed5ef987a862a284cf024db50ce3f73b4545c7e36800ae1ad94759248c140a000000808e73a44475925ae478b9e123040373fddf5f5b7475d8392dfe08c8c0e0bc1d2bea4e7c8b0414fed471444a17421c80eb0390eb8f43097557fc801fbfe53a25ff672da011f7b419fb8ac6176f75133aae830c5f3354114d9a3be0ecc215dadb407d5730a8cf005683ebb06f3249edbe04eb9870ed7fe7a680e477ac4b799ba0ee";

  const deploy = async () => {
    acc = await AccessControl.new(pubKey, {
      from: owner
    });
  };

  before(deploy);

  describe("TestGeneralProps", async () => {
    it("should verify the contract owner", async () => {
      retrievedOwner = await acc.owner();
      assert.equal(retrievedOwner, owner);
    });
    it("should verify the pubKey", async () => {
      retrievedPubKey = await acc.pubKey();
      assert.equal(retrievedPubKey, pubKey);
    });
  });

  describe("TestGrantAccess", async () => {
    it("should check the request access method", async () => {
      let result = await acc.requestAccess({
        from: requester
      });

      truffleAssert.eventEmitted(
        result,
        "AccessRequested",

        ev => {
          assert.equal(ev[0], requester);
          return ev[0] === requester;
        }
      );

      let acl = await acc.acl(requester);
      assert.equal(true, acl.pending);
      assert.equal("", acl.key);
    });

    it("should check the grant access method", async () => {
      //NOTE public key can somehow not be recovered when using Ganache (yields wrong address)
      // const signer = ethCrypto.recoverPublicKey(
      //   "0xc04b809d8f33c46ff80c44ba58e866ff0d5..", // signature
      //   EthCrypto.hash.keccak256("foobar") // message hash
      // );

      const publicKey = ethCrypto.publicKeyByPrivateKey(requesterPrivKey);
      let addr = ethCrypto.publicKey.toAddress(publicKey);
      assert.equal(addr, requester);

      const encrypted = await ethCrypto.encryptWithPublicKey(
        publicKey, // publicKey
        prvKey // message
      );

      const encryptedKey = ethCrypto.cipher.stringify(encrypted);

      let attributes = ["admin", "it_departement"];
      let attributesBytes = [];
      attributes.forEach(attr => {
        let attrHex = web3.utils.asciiToHex(attr);
        let attrBytes = web3.utils.hexToBytes(attrHex);
        attributesBytes.push(attrBytes);
      });

      let result = await acc.grantAccess(
        requester,
        encryptedKey,
        attributesBytes,
        {
          from: owner
        }
      );

      truffleAssert.eventEmitted(result, "AccessGranted", ev => {
        assert.equal(ev[0], requester);
        assert.equal(ev[1], encryptedKey);

        let hexAttrs = [];
        attributesBytes.forEach(attr => {
          let hexAttr = web3.utils.bytesToHex(attr);
          hexAttr = web3.utils.padRight(hexAttr, 64); //why 64?
          hexAttrs.push(hexAttr);
        });

        assert.deepStrictEqual(ev[2], hexAttrs);

        return true;
      });

      const acl = await acc.acl(requester);
      const parsedEncryptedPrivKey = ethCrypto.cipher.parse(acl.key);
      const decryptedKey = await ethCrypto.decryptWithPrivateKey(
        requesterPrivKey,
        parsedEncryptedPrivKey
      );
      assert.equal(decryptedKey, prvKey);
    });

    it("should check that the grant access fails for other user than owner", async () => {
      let attributes = ["admin", "it_departement"];
      let attributesBytes = [];
      attributes.forEach(attr => {
        let attrHex = web3.utils.asciiToHex(attr);
        let attrBytes = web3.utils.hexToBytes(attrHex);
        attributesBytes.push(attrBytes);
      });

      await truffleAssert.reverts(
        acc.grantAccess(requester, prvKey, attributesBytes, {
          from: requester
        }),
        "Caller is not owner"
      );
    });
  });

  describe("TestDevicePolicies", async () => {
    it("should add a device policy", async () => {
      const policy = "(admin AND date > 15434456322)";
      let result = await acc.setDevicePolicy(iotDevice, policy, {
        from: owner
      });

      truffleAssert.eventEmitted(result, "DevicePolicyUpdated", ev => {
        assert.equal(ev[0], iotDevice);
        assert.equal(ev[1], policy);
        return true;
      });

      let retrievedPolicy = await acc.devices(iotDevice);
      assert.equal(retrievedPolicy, policy);
    });
    it("should not add a device policy from other user than owner", async () => {
      const policy = "(admin AND date > 15434456322)";

      await truffleAssert.reverts(
        acc.setDevicePolicy(iotDevice, policy, {
          from: requester
        }),
        "Caller is not owner"
      );
    });
    it("should remove a device policy", async () => {
      let result = await acc.removeDevicePolicy(iotDevice, {
        from: owner
      });

      truffleAssert.eventEmitted(result, "DevicePolicyDeleted", ev => {
        assert.equal(ev[0], iotDevice);
        return true;
      });

      let retrievedPolicy = await acc.devices(iotDevice);
      assert.equal(retrievedPolicy, "");
    });
    it("should not remove a device policy if not owner", async () => {
      await truffleAssert.reverts(
        acc.removeDevicePolicy(iotDevice, {
          from: requester
        }),
        "Caller is not owner"
      );
    });
  });
});
