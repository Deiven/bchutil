// Copyright (c) 2014-2016 The btcsuite developers
// Copyright (c) 2018 The gcash developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package merkleblock_test

import (
	"bytes"
	"encoding/hex"
	"testing"

	"github.com/gcash/bchd/chaincfg/chainhash"
	"github.com/gcash/bchd/wire"
	"github.com/gcash/bchutil"
	"github.com/gcash/bchutil/bloom"
	"github.com/gcash/bchutil/merkleblock"
)

// TestMerkleBlock3 tests merkleblock encoding using bloom filter. This test
// derives from when merkleblock encoding functionality was bundled under
// the bloom filter package as TestMerkleBlock3
func TestNewMerkleBlockWithFilter(t *testing.T) {
	blockStr := "0100000079cda856b143d9db2c1caff01d1aecc8630d30625d10e8b" +
		"4b8b0000000000000b50cc069d6a3e33e3ff84a5c41d9d3febe7c770fdc" +
		"c96b2c3ff60abe184f196367291b4d4c86041b8fa45d630101000000010" +
		"00000000000000000000000000000000000000000000000000000000000" +
		"0000ffffffff08044c86041b020a02ffffffff0100f2052a01000000434" +
		"104ecd3229b0571c3be876feaac0442a9f13c5a572742927af1dc623353" +
		"ecf8c202225f64868137a18cdd85cbbb4c74fbccfd4f49639cf1bdc94a5" +
		"672bb15ad5d4cac00000000"
	blockBytes, err := hex.DecodeString(blockStr)
	if err != nil {
		t.Errorf("TestMerkleBlock3 DecodeString failed: %v", err)
		return
	}
	blk, err := bchutil.NewBlockFromBytes(blockBytes)
	if err != nil {
		t.Errorf("TestMerkleBlock3 NewBlockFromBytes failed: %v", err)
		return
	}

	f := bloom.NewFilter(10, 0, 0.000001, wire.BloomUpdateAll)

	inputStr := "63194f18be0af63f2c6bc9dc0f777cbefed3d9415c4af83f3ee3a3d669c00cb5"
	hash, err := chainhash.NewHashFromStr(inputStr)
	if err != nil {
		t.Errorf("TestMerkleBlock3 NewHashFromStr failed: %v", err)
		return
	}

	f.AddHash(hash)

	mBlock, _ := merkleblock.NewMerkleBlockWithFilter(blk, f)

	wantStr := "0100000079cda856b143d9db2c1caff01d1aecc8630d30625d10e8b4" +
		"b8b0000000000000b50cc069d6a3e33e3ff84a5c41d9d3febe7c770fdcc" +
		"96b2c3ff60abe184f196367291b4d4c86041b8fa45d630100000001b50c" +
		"c069d6a3e33e3ff84a5c41d9d3febe7c770fdcc96b2c3ff60abe184f196" +
		"30101"
	want, err := hex.DecodeString(wantStr)
	if err != nil {
		t.Errorf("TestMerkleBlock3 DecodeString failed: %v", err)
		return
	}

	got := bytes.NewBuffer(nil)
	err = mBlock.BchEncode(got, wire.ProtocolVersion, wire.LatestEncoding)
	if err != nil {
		t.Errorf("TestMerkleBlock3 BchEncode failed: %v", err)
		return
	}

	if !bytes.Equal(want, got.Bytes()) {
		t.Errorf("TestMerkleBlock3 failed merkle proof comparison: "+
			"got %v want %v", got.Bytes(), want)
		return
	}
}

// TestValidNewMerkleBlockWithTxnSet tests encoding a merkle proof for a given block
// with specified transaction set
func TestValidNewMerkleBlockWithTxnSet(t *testing.T) {

	tests := []struct {
		name      string
		blockHash string
		// the fully encoded raw block
		blockStr string
		// the list of transactions to be merkle proof for
		txnIds []*chainhash.Hash
		proof  string
	}{
		{
			name:      "Generate 2 transaction proof for Testnet block 1253848",
			blockHash: "00000000000002999c8b785403bd0a9f262e2c2ec719965105aff243d3343eeb",
			blockStr:  "00000020c2981857b4516c746e24199820dd2309818a058ce5371a27be0000000000000091bfb84d0fce1e261d97c86d84c44ff66aab3a610e68ac1cf09873159ed1ce9a56ae825b1013041a21aebff30602000000010000000000000000000000000000000000000000000000000000000000000000ffffffff2e03d821130456ae825b68747470733a2f2f6769746875622e636f6d2f62636578742f010000147b60020000000000ffffffff01c3a95009000000001976a91448b20e254c0677e760bab964aec16818d6b7134a88ac000000000100000001814ec7ade395fe51aaf3495856e0ef9406e4dda8bc3f6288cde90d4030e1e9b1010000006a4730440220518f58ff653bd35ec95f8ddee1de76ecde8afd69e95d0a7f24cceef9ebb58392022058d810cd26ffa5d5bfc2e7f07219fe2f2aaf0ce711d1348a15da9ba00815085c412102e76a7019e3b5ecc368f0ae16573216132575db3baa137d80c423205ec41c36b4ffffffff025a250200000000001976a9145210cff0fc86b990a84a994f3508770cb6de592788ac0d660800000000001976a914ba90188f217a59720821ba42db80ee7fda0aaf8988ac00000000020000000193f40f40cd6aaf9c80781911facd3732d9b94f0fb752e714e5ae71774c2a6492000000006a47304402201462d9d9585483696e267c050ffef07331e45cacd476f3c6f342b7a79094013602202465bebe15dbd93f7de145acfc17883e750f0511501bdd65643553d558adcc0e412103295c227557eeb4c973ce576ac6d38faad3c78cbec41e260c4a5a09713699d9e9ffffffff02805c240e000000001976a91411838d8f4aaafd5bf45494f02ee0b6e8eeba212488ac0000000000000000336a3143704d3115b27e9b1366c275dd15ea0950173e160194d3c794f5f61f8521d215e90f773f6f37651400ea62011b065fc25300000000010000000167028a425578d0391deadc8724e0b60c49c4ac196b62de8907519625819d9af7000000006b483045022100d4252d98872f9a19a0bb4b2b9fd80d3839eed173a99c3e250c2cd8460d2c47f9022007a9e9ec47e343b17c47d6c0583fa3dccdfde8a326603f25be105f8b3167a9f7412102ddc55a123521a0f4c371dc848338e60237efc4c2d322693414f3978066412485ffffffff0378ae5f01000000001976a9146add93b82bb1f549e5ef68c255569886fd64d4c488ac7c080f00000000001976a914f177271bf2d4169fab6602170105bf4282d07e1f88acbf380000000000001976a914d56e43122efdc36869464336805f72bd07514a4e88ac0000000001000000014846033a92b7038518db97dc6a522c7f386084c15109001413daf6ceb4ba4d8f000000006b483045022100e54aa3f8854d139c7661f6248e114a8e9353d4e3c172c258ccbeef6d0d05b0fd02201162c6bf0a140099d320d36c6800dcf15bbd3270d75eb9df5a6ffb663dc0f1ba412102519a0a15916fdb8f949864193a7deed81a0c06c46e93afa811589811a29bf417ffffffff03386c5001000000001976a9142b05471b6bd4cfabb88e92015dc9442bc9e397c488aca8070f00000000001976a914e192fc7af84a35789e5946010aa3277417cfe9b888ac93390000000000001976a914f6221d2df9b48cf5cf4454bd2d5d677f869da5f088ac000000000100000001f9094fcfef8ec0a7024f70c415af1510fec18db826bf914695871b6b4c5a1e78000000006b483045022100a6808939541f454a875916e188996dca462e48f9a91c65516a2f88c90197be2602202de82a3481272dd10824b018a402f97eed22593e1e40385b34fd4174740539ba412102b1f8b43a4192a8524509b5aa95f398ba23f646b74b8b7bcfbdc1b48461148500ffffffff03b8d5b700000000001976a9144f71a322d45ea14ce1bb89823100b2d3aa35efdc88ac1ee69300000000001976a91411936a43a4f0401c94205e4a303bcc85b85fa4b588ac5daf0400000000001976a914f1d273facea139374c1d2b52050fab5af637c7f088ac00000000",
			txnIds: []*chainhash.Hash{
				hashFromStr("bfc591e1471b1f247258b5274a845e666e5bc9b918a154e1cf4c6211b640711d"),
				hashFromStr("8f4dbab4cef6da1314000951c18460387f2c526adc97db188503b7923a034648"),
			},
			proof: "00000020c2981857b4516c746e24199820dd2309818a058ce5371a27be0000000000000091bfb84d0fce1e261d97c86d84c44ff66aab3a610e68ac1cf09873159ed1ce9a56ae825b1013041a21aebff3060000000411303ed124cfa04609d6728e78ae56ac913a166da35ae1338d268a05004f82741d7140b611624ccfe154a118b9c95b6e665e844a27b55872241f1b47e191c5bf4846033a92b7038518db97dc6a522c7f386084c15109001413daf6ceb4ba4d8fea3c13cb09d08cfa3e43a6a40d3d01602835093c48d414bcce0777267b8190fc013b",
		},
		{
			name:      "Generate 1 transaction proof for Testnet block 1268825",
			blockHash: "0000000000000012091a94079b64c6ecf92df6744d973b18e8c5d89e7b45d66b",
			blockStr:  "00e0002075dbd04988a32f7fe6346a7908f04d0aca6f3cad22c6c138fe020000000000000679e44eddbcd820b33c4a287f66df403163576c02c05be6bb509aaca224dfc68ff7f25be142031ad3c640580201000000010000000000000000000000000000000000000000000000000000000000000000ffffffff2703595c131a4d696e656420627920416e74506f6f6c20a7000b01205bf2f76401000000bf080000ffffffff02002ba804000000001976a914d1746ca2a02768a93abf5636fdec76b782f3b40c88ac0000000000000000266a24aa21a9ed91250cd7d5cf941e75467ceb21af3197cdd138610cb49c5ee467f594c63af16a0000000001000000017bca6219086abcf83eed7a8030d8d9c26adba59801d82e318493db96c21bef70010000006b483045022100f9f3265233ae9ebafe42c04f6a6a03be4473cf0558d37fee241cede72ec862fb02200b2e3f6f09b54fa52d90916aa176626e80bcfaf00d5057e5cc3951f89a5e4c5d4121024b18442b84450bd3040cfc03646e1d773f35f347df5e01a530ceef81d3cb44c7ffffffff0200000000000000002c6a09696e7465726c696e6b20f6384b8d6875acf4d9ff3e14311cc46761066baf95f6310261eaa860fee3d86ab8419d04000000001976a914f56c4f877ae2d24c61333629e9a9c833b005e0ef88ac00000000",
			txnIds: []*chainhash.Hash{
				hashFromStr("650394f9753a3c2b138ef1ee5d98feebc7c37fe83464dc2b83773c1e6808316b"),
			},
			proof: "00e0002075dbd04988a32f7fe6346a7908f04d0aca6f3cad22c6c138fe020000000000000679e44eddbcd820b33c4a287f66df403163576c02c05be6bb509aaca224dfc68ff7f25be142031ad3c64058020000000269430520d964b07d7ee8724f8a7c2fe391c38cef14d0c5c9ad5030f9bb04c4bc6b3108681e3c77832bdc6434e87fc3c7ebfe985deef18e132b3c3a75f99403650105",
		},
		{
			name:      "Generate 8 transaction proof for Testnet block 1268830",
			blockHash: "0000000000030097fccee7049aa3ae6b8aa26abd2fa556965ba24106b811770f",
			blockStr:  "000000209fb15bb4903b1b2ca439b082e76d05ffb4edc95ffaf9864e73010000000000007c8ffd4aa7d4bb1ac7021d9823b4d4a11e11c7a404a7185d898c43eabebc8ecdaa04f35bffff001d36d97a580b02000000010000000000000000000000000000000000000000000000000000000000000000ffffffff25035e5c1304aa04f35b2f44504f4f4c2e544f502f544553542f01000006000000081f0c0000ffffffff01ec0eaa04000000001976a914aece60e928163f56e10d4f048acab182f8f3697288ac000000000200000001304f7c2ad9d17b34a968db54ccfa9c4e1b025947e4a82c40321b19fed8779f8b000000006a47304402200b388ae2b6854c1c67a307ee1c8b954a133a04de6434b1f80e9769a70784b12702203c77a07af4f14c8e3360f23681bdc3c8e3a4739140d4767bd75425fb3b2d73b14121020d01a18f138e0461e14fa3bac7913be2ad9e26b9a9fc4cfded431344ca2a0dcbffffffff027da2760b000000001976a91448374814c1f3c3138e2000d5d5e04bdedb25cf5088aca0860100000000001976a9144073d94cf154f4670c87b1ed17f30851b190084f88ac000000000100000001da62ece2344d78f39074bffae837a280d3032e4db49308162895ad3fead77dda000000006a4730440220797ecb34d7802d1a46bc8ec69401ea3d0e4fab545e18d3ec43f7939f21a14260022063ed13285c258962f9e2f3706beaa1084f1b3343c936df9a70cd046e2cac99fa41210290e0c3a165df07b83a6e93051236afcf4e19a6e5f4439ddc07ffd4fc74742896feffffff0200000000000000002c6a09696e7465726c696e6b20c447ebefe381ae70be47f730a55f46ab74e1c9fce742c1dc22dd098ff32592ddb07bd705000000001976a914876b90552d174a69e46c8ad25f0d6b09e165031c88ac000000000200000001ef1d12cc76a92d1491dd8e31574ff5d780b7ff0740196bd0efdf49157f5d51b8010000006b4830450221009c926a8bb7b5de293ce4c674848682628380bd68ac8c3c26246b820b87b014aa022063d69e646335e02159cacdde060a29a0edeabe900b9902177ef76bbef1e917c3412102426c37d9d08afb7885a76b06ab4689d457aaccaba34fb5524a9c6af6dea55e19feffffff024b13cd000000000017a9145a0dbd8222e4e96629d5da325877bc95aa42922787ad4cb714190000001976a91421b97f5dd8d48c0f2c05164d5bc5510113ab60df88ac5d5c130001000000014f3721caa87e688ec476932c546b3e33c65698e99250c2fdff6bdc18fed3efe9000000006b4830450221008d71c1125f5f6d4d73f9762b6e2df380cfa7db602c46f84e76a928eb9c39dc3e022057cf7514b6476f33691c3b42b364c4014ae938ac94b5b1279d42dea7d537a069412102e01ca9378f18e3dbbc5fe71e51e85a3e0d169ef22d19caa6f126fbb1bdf67f98feffffff0200000000000000002c6a09696e7465726c696e6b20aa639003d31ba42fe64553545bc8d0ad88ed0a3c2027fa2fb92d10486b6b9499ae7dd705000000001976a914874c91285a2adbe296b4c1686322cca7ad21b43f88ac000000000200000002bc35365adc9258e44055fe882156ff3f577cb298ae2ef95a6a15edf1fb5059ce010000006b48304502210092ec4652adbf5a2fb4521cb361dfce6b3353679428b073d07d0b54ebfc18983e0220688c5c3593db3a6b3907ff6a07aef2f857ecdb811af0bd5c68dbaa85e866678641210279d6ea7c43b8cec027b91c36b59d9a9fb1dca501e3db5286f6cc1ea5e10f0370feffffff14bc71273d1f73ebdbc83e8f077bbdf13d21d10160c72a4ee76bb7681eecd50c010000006b483045022100ccc3b745aa2e3f4dc8dd479fee459cf2e81585e38a77b52fa3da34611119b60a02207ee767e40b56a8aca37d8650146a5b920c3743e4df51693430390b5fd86214ca412103987ff61fbb7f041f0c936c091c13f32236a9211e87569c1e3b59780964f943c4feffffff023418d704000000001976a914d61b789c70eae758a4125dd0ecd012a2fc16d17588ac6ce54200000000001976a9148b2621ce7a6ed1a567cbe97e04642194a491e53788ac3d5c13000200000001bc35365adc9258e44055fe882156ff3f577cb298ae2ef95a6a15edf1fb5059ce000000006b483045022100bcd762a88de9a582690980fc098bc5dc9651f0ec66543419a87a1ef819d93096022053ca90a8c5ad1e9f4cc724e6b04b584f278b24d644b7f622744cc4a61f9085ea41210307cce48adeab39e7dd64964789807b07eb15f84d7e4e1542ab67d7cb597f22affeffffff02aa654b00000000001976a9149743538ab2599f34bf38753a6b7e4463de7d109c88ac1caef201000000001976a914c5ea622b98e9c505e94c26c562b7f59fd2a3e0d788ac5a5c13000100000001db9ad3ec9f15c745b3c90891848dcf649bd47f4f354a3f48537d319f33effec5010000006a47304402204e65d509e81e3e754a0ee0e7fa527028727018f7228b7829f85ed0a447dd135a022020ffaf80e7519b388febc1ecb4c75b85e10161aa2273eb2146112a0c12f52b4841210380081eb871d78f11bced442262bc8b57eb5ad239a6f2fb383608219c374ae665ffffffff0200000000000000002c6a09696e7465726c696e6b20c447ebefe381ae70be47f730a55f46ab74e1c9fce742c1dc22dd098ff32592dda0e19c04000000001976a914de3ffa1dab929f0f06eaa416ce58ed6d2e3d103288ac000000000100000001ac2c38d41688322cb1f87c209f9a40ed696d86c46f6e27ac2c203f552a38b805010000006a473044022004b186dd48ddb9554fb6cc44970ce11a1aa74f786c1467923f16f6b92a10cd390220415e14a1ab322fa2944b1c9c247d19784bb376204ca10c4d3dcd11e27c3d2b0d412103eb52cee1016814e29d28ec6177bfd6dc301181336425e07628d4bd406b67aaceffffffff0250c300000000000017a9149766fed17ed78584f14872dfe303cf2c49e7131187de22ec00000000001976a914f330a4f9d1ce9a105a20ba124c023cf24afd808a88ac0000000001000000012b9b225d7acd0a297f5300c311fb734f17258e65464338e113ff850245da2d41010000006a473044022050fb8862fa1b2e37c84beb91baefd3c9f866f717ddfbc7c271fbb961cbc2f246022037d62676a66a09c9404eabd6b3ed491c66ceace9a67c18207c86c9421991d223412103136a352ddcfe27c87bd6ac000cb4c0cbb12f6622c69dd953db465d40d7642ff9feffffff02af7cd705000000001976a91411d1669803fc39e45b4b77d42dd70bdd60e88d2e88ac00000000000000002c6a09696e7465726c696e6b200d0ca64ccb70e1fb9b4286a12403594df749f52bd1fcb5491b800e801a84a7120000000001000000014634d50fcc3eeba3aaa5d1ab53d9a86bb25400ffc62294bab99653c80ed19225010000006b483045022100ef716b93028a9b0f479d4ab7d5adac078af6a34cf72cf04e7983969bc0049bad0220090fa9c27c5f173f390b5c6dbb430c50eb8cff9c8d0a8989e4c8dc61737cd38141210259d50e231bd5ba5ef874595a0b2ec52e5f69500e6f7aa2dbb067038469d7572bfeffffff02ad7ed705000000001976a914f3e8e35a7ea2be689c8a7425bf221c47bc935d1688ac00000000000000002c6a09696e7465726c696e6b20be1004860213cfa4dc7a6d1fed22fd6e29b6f7f6a30438178bc9cf16f1e6997d00000000",
			txnIds: []*chainhash.Hash{
				hashFromStr("c9b5f17d01104d35f7daa5010aec9bdf9ba6717ef8ab472e937ecc73b019a83b"),
				hashFromStr("00c412161f5d61796048eff1b9cc6976395023578c96393f74215df863931aa5"),
				hashFromStr("0ea54593df07295cff857103163ae3d8c5ad80adadd5d8347ebd7a278bef33e7"),
				hashFromStr("412dda450285ff13e1384346658e25174f73fb11c300537f290acd7a5d229b2b"),
				hashFromStr("506f9c2e56a0de8a84cab28da5be434a4ce8affd66d02429a7f39c61d7dbe989"),
				hashFromStr("aeb1f9d6885aedf2da78e39bee2df14cd456ad6dd0474c6f1048b819cd1ec62b"),
				hashFromStr("b66c92b01cd81737a4ddf5888fe3366f1cd2ddeea1c523b1610d402735de6eb6"),
				hashFromStr("e9efd3fe18dc6bfffdc25092e99856c6333e6b542c9376c48e687ea8ca21374f"),
			},
			proof: "000000209fb15bb4903b1b2ca439b082e76d05ffb4edc95ffaf9864e73010000000000007c8ffd4aa7d4bb1ac7021d9823b4d4a11e11c7a404a7185d898c43eabebc8ecdaa04f35bffff001d36d97a580b0000000b3ba819b073cc7e932e47abf87e71a69bdf9bec0a01a5daf7354d10017df1b5c9a51a9363f85d21743f39968c572350397669ccb9f1ef486079615d1f1612c400e733ef8b277abd7e34d8d5adad80adc5d8e33a16037185ff5c2907df9345a50e4ce18fcbcc42286439c855400b5fb8f8969ffa185ee7a34ae76bebd36b26163e2b9b225d7acd0a297f5300c311fb734f17258e65464338e113ff850245da2d4189e9dbd7619cf3a72924d066fdafe84c4a43bea58db2ca848adea0562e9c6f500e87603aeca0ebf464ab72f7d1d1743218c121ac734cedced590ae515c0d27882bc61ecd19b848106f4c47d06dad56d44cf12dee9be378daf2ed5a88d6f9b1aeb66ede3527400d61b123c5a1eeddd21c6f36e38f88f5dda43717d81cb0926cb6da62ece2344d78f39074bffae837a280d3032e4db49308162895ad3fead77dda4f3721caa87e688ec476932c546b3e33c65698e99250c2fdff6bdc18fed3efe903ffbe6f",
		},
	}

	for _, tt := range tests {
		tt := tt // capture range variable
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// decode block string
			blockBytes, err := hex.DecodeString(tt.blockStr)

			if err != nil {
				t.Errorf("blockBytes DecodeString failed: %v", err)
				return
			}

			blk, err := bchutil.NewBlockFromBytes(blockBytes)

			if err != nil {
				t.Errorf("NewBlockFromBytes failed: %v", err)
				return
			}

			// create merkle proof
			mBlock, matchedIndices := merkleblock.NewMerkleBlockWithTxnSet(blk, tt.txnIds)

			// check number of indices matches number of transactions given
			if len(matchedIndices) != len(tt.txnIds) {
				t.Errorf("indices mismatch, got %d wanted %d", len(matchedIndices), len(tt.txnIds))
				return
			}

			// encode wire message
			got := bytes.NewBuffer(nil)
			err = mBlock.BchEncode(got, wire.ProtocolVersion, wire.LatestEncoding)

			if err != nil {
				t.Errorf("mBlock BchEncode failed: %v", err)
				return
			}

			want, err := hex.DecodeString(tt.proof)

			if err != nil {
				t.Errorf("want Proof DecodeString failed: %v", err)
				return
			}

			// compare result
			if !bytes.Equal(want, got.Bytes()) {
				t.Errorf("Failed merkle proof comparison: "+
					"got %v want %v", got.Bytes(), want)
			}

		})
	}
}

func TestInvalidNewMerkleBlockWithTxnSet(t *testing.T) {

	tests := []struct {
		name      string
		blockHash string
		// the fully encoded raw block
		blockStr string
		// the list of transactions to be merkle proof for
		txnIds  []*chainhash.Hash
		indices int
	}{
		{
			name:      "Transactions from different block",
			blockHash: "0000000000000012091a94079b64c6ecf92df6744d973b18e8c5d89e7b45d66b",
			blockStr:  "00e0002075dbd04988a32f7fe6346a7908f04d0aca6f3cad22c6c138fe020000000000000679e44eddbcd820b33c4a287f66df403163576c02c05be6bb509aaca224dfc68ff7f25be142031ad3c640580201000000010000000000000000000000000000000000000000000000000000000000000000ffffffff2703595c131a4d696e656420627920416e74506f6f6c20a7000b01205bf2f76401000000bf080000ffffffff02002ba804000000001976a914d1746ca2a02768a93abf5636fdec76b782f3b40c88ac0000000000000000266a24aa21a9ed91250cd7d5cf941e75467ceb21af3197cdd138610cb49c5ee467f594c63af16a0000000001000000017bca6219086abcf83eed7a8030d8d9c26adba59801d82e318493db96c21bef70010000006b483045022100f9f3265233ae9ebafe42c04f6a6a03be4473cf0558d37fee241cede72ec862fb02200b2e3f6f09b54fa52d90916aa176626e80bcfaf00d5057e5cc3951f89a5e4c5d4121024b18442b84450bd3040cfc03646e1d773f35f347df5e01a530ceef81d3cb44c7ffffffff0200000000000000002c6a09696e7465726c696e6b20f6384b8d6875acf4d9ff3e14311cc46761066baf95f6310261eaa860fee3d86ab8419d04000000001976a914f56c4f877ae2d24c61333629e9a9c833b005e0ef88ac00000000",
			txnIds: []*chainhash.Hash{
				// valid txn from block height 1268825
				hashFromStr("650394f9753a3c2b138ef1ee5d98feebc7c37fe83464dc2b83773c1e6808316b"),
				// invalid txn from block height 1268830
				hashFromStr("b66c92b01cd81737a4ddf5888fe3366f1cd2ddeea1c523b1610d402735de6eb6"),
			},
			indices: 1,
		},
	}

	for _, tt := range tests {
		tt := tt // capture range variable
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// decode block string
			blockBytes, err := hex.DecodeString(tt.blockStr)

			if err != nil {
				t.Errorf("blockBytes DecodeString failed: %v", err)
				return
			}

			blk, err := bchutil.NewBlockFromBytes(blockBytes)

			if err != nil {
				t.Errorf("NewBlockFromBytes failed: %v", err)
				return
			}

			// create merkle proof
			_, matchedIndices := merkleblock.NewMerkleBlockWithTxnSet(blk, tt.txnIds)

			if len(matchedIndices) != tt.indices {
				t.Errorf("Indices mismatch, got %d expected %d", len(matchedIndices), tt.indices)
			}

		})
	}
}
