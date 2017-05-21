package multihash

import (
	"bytes"
	"encoding/hex"
	"testing"
)

type SumTestCase struct {
	code   uint64
	length int
	input  string
	hex    string
}

var sumTestCases = []SumTestCase{
	SumTestCase{ID, 3, "foo", "0003666f6f"},
	SumTestCase{SHA1, -1, "foo", "11140beec7b5ea3f0fdbc95d0dd47f3c5bc275da8a33"},
	SumTestCase{SHA1, 10, "foo", "110a0beec7b5ea3f0fdbc95d"},
	SumTestCase{SHA2_256, -1, "foo", "12202c26b46b68ffc68ff99b453c1d30413413422d706483bfa0f98a5e886266e7ae"},
	SumTestCase{SHA2_256, 16, "foo", "12102c26b46b68ffc68ff99b453c1d304134"},
	SumTestCase{SHA2_512, -1, "foo", "1340f7fbba6e0636f890e56fbbf3283e524c6fa3204ae298382d624741d0dc6638326e282c41be5e4254d8820772c5518a2c5a8c0c7f7eda19594a7eb539453e1ed7"},
	SumTestCase{SHA2_512, 32, "foo", "1320f7fbba6e0636f890e56fbbf3283e524c6fa3204ae298382d624741d0dc663832"},
	SumTestCase{SHA3, 32, "foo", "14204bca2b137edc580fe50a88983ef860ebaca36c857b1f492839d6d7392452a63c"},
	SumTestCase{SHA3_512, 16, "foo", "14104bca2b137edc580fe50a88983ef860eb"},
	SumTestCase{SHA3_512, -1, "foo", "14404bca2b137edc580fe50a88983ef860ebaca36c857b1f492839d6d7392452a63c82cbebc68e3b70a2a1480b4bb5d437a7cba6ecf9d89f9ff3ccd14cd6146ea7e7"},
	SumTestCase{SHA3_224, -1, "beep boop", "171c0da73a89549018df311c0a63250e008f7be357f93ba4e582aaea32b8"},
	SumTestCase{SHA3_224, 16, "beep boop", "17100da73a89549018df311c0a63250e008f"},
	SumTestCase{SHA3_256, -1, "beep boop", "1620828705da60284b39de02e3599d1f39e6c1df001f5dbf63c9ec2d2c91a95a427f"},
	SumTestCase{SHA3_256, 16, "beep boop", "1610828705da60284b39de02e3599d1f39e6"},
	SumTestCase{SHA3_384, -1, "beep boop", "153075a9cff1bcfbe8a7025aa225dd558fb002769d4bf3b67d2aaf180459172208bea989804aefccf060b583e629e5f41e8d"},
	SumTestCase{SHA3_384, 16, "beep boop", "151075a9cff1bcfbe8a7025aa225dd558fb0"},
	SumTestCase{DBL_SHA2_256, 32, "foo", "5620c7ade88fc7a21498a6a5e5c385e1f68bed822b72aa63c4a9a48a02c2466ee29e"},
	SumTestCase{BLAKE2B_MAX, 64, "foo", "c0e40240ca002330e69d3e6b84a46a56a6533fd79d51d97a3bb7cad6c2ff43b354185d6dc1e723fb3db4ae0737e120378424c714bb982d9dc5bbd7a0ab318240ddd18f8d"},
	SumTestCase{BLAKE2B_MAX - 32, 32, "foo", "a0e40220b8fe9f7f6255a6fa08f668ab632a8d081ad87983c77cd274e48ce450f0b349fd"},
	SumTestCase{BLAKE2B_MAX - 16, 32, "foo", "b0e40220e629ee880953d32c8877e479e3b4cb0a4c9d5805e2b34c675b5a5863c4ad7d64"},
	SumTestCase{BLAKE2S_MAX, 32, "foo", "e0e4022008d6cad88075de8f192db097573d0e829411cd91eb6ec65e8fc16c017edfdb74"},
	SumTestCase{MURMUR3, 4, "beep boop", "2204243ddb9e"},
	SumTestCase{KECCAK_224, -1, "beep boop", "1a1c2bd72cde2f75e523512999eb7639f17b699efe29bec342f5a0270896"},
	SumTestCase{KECCAK_256, 32, "foo", "1b2041b1a0649752af1b28b3dc29a1556eee781e4a4c3a1f7f53f90fa834de098c4d"},
	SumTestCase{KECCAK_384, -1, "beep boop", "1c300e2fcca40e861fc425a2503a65f4a4befab7be7f193e57654ca3713e85262b035e54d5ade93f9632b810ab88b04f7d84"},
	SumTestCase{KECCAK_512, -1, "beep boop", "1d40e161c54798f78eba3404ac5e7e12d27555b7b810e7fd0db3f25ffa0c785c438331b0fbb6156215f69edf403c642e5280f4521da9bd767296ec81f05100852e78"},
	SumTestCase{SHAKE_128, 32, "foo", "1820f84e95cb5fbd2038863ab27d3cdeac295ad2d4ab96ad1f4b070c0bf36078ef08"},
	SumTestCase{SHAKE_256, 64, "foo", "19401af97f7818a28edfdfce5ec66dbdc7e871813816d7d585fe1f12475ded5b6502b7723b74e2ee36f2651a10a8eaca72aa9148c3c761aaceac8f6d6cc64381ed39"},
	SumTestCase{SKEIN256_MAX - 31, -1, "beep boop", "2c"},
	SumTestCase{SKEIN256_MAX - 30, -1, "beep boop", "7025"},
	SumTestCase{SKEIN256_MAX - 29, -1, "beep boop", "096729"},
	SumTestCase{SKEIN256_MAX - 28, -1, "beep boop", "ec2226c1"},
	SumTestCase{SKEIN256_MAX - 27, -1, "beep boop", "cf8090aee1"},
	SumTestCase{SKEIN256_MAX - 26, -1, "beep boop", "b7450f6303fb"},
	SumTestCase{SKEIN256_MAX - 25, -1, "beep boop", "c621478d4c75ea"},
	SumTestCase{SKEIN256_MAX - 24, -1, "beep boop", "9a1188a5e2667bb6"},
	SumTestCase{SKEIN256_MAX - 23, -1, "beep boop", "faa0212e6d42b4c28e"},
	SumTestCase{SKEIN256_MAX - 22, -1, "beep boop", "4055fa5dc32904feea71"},
	SumTestCase{SKEIN256_MAX - 21, -1, "beep boop", "df8f80b280d62e0979fdb5"},
	SumTestCase{SKEIN256_MAX - 20, -1, "beep boop", "b0d9c7f6e423d051be049dc9"},
	SumTestCase{SKEIN256_MAX - 19, -1, "beep boop", "4a9ee56765345ff5f2f5042fa1"},
	SumTestCase{SKEIN256_MAX - 18, -1, "beep boop", "511ce50e2c7fd7323d9d25c9d790"},
	SumTestCase{SKEIN256_MAX - 17, -1, "beep boop", "df2d0eb8ef39a5e6f4e4919182f22f"},
	SumTestCase{SKEIN256_MAX - 16, -1, "beep boop", "74392144ed4ea9ec729d789ba668febf"},
	SumTestCase{SKEIN256_MAX - 15, -1, "beep boop", "45449099727f369a66920f6a5316d6088b"},
	SumTestCase{SKEIN256_MAX - 14, -1, "beep boop", "b8cee8b0aa18a1e0ed08f6297639a5bb2bc7"},
	SumTestCase{SKEIN256_MAX - 13, -1, "beep boop", "0a5794defb06823c67910cc5bc7d2cacb448b7"},
	SumTestCase{SKEIN256_MAX - 12, -1, "beep boop", "d09f39744f318d689b90971f535793cb821e47c7"},
	SumTestCase{SKEIN256_MAX - 11, -1, "beep boop", "8500fdbf961737b9df2ada6e15bf7385c1c9c89cf6"},
	SumTestCase{SKEIN256_MAX - 10, -1, "beep boop", "127d5e617621dc09dd75b0010b87cca5ce1f8373f3fa"},
	SumTestCase{SKEIN256_MAX - 9, -1, "beep boop", "513844f6e58157c22c5224ce4f008524424d339711e6d1"},
	SumTestCase{SKEIN256_MAX - 8, -1, "beep boop", "e6a15dc094eb2f0cdb76c974383b79dc55b6a2d92c5dae7f"},
	SumTestCase{SKEIN256_MAX - 7, -1, "beep boop", "4829b11bf52229fb5c8e673194daa1abf0284b5f3ce9bb292b"},
	SumTestCase{SKEIN256_MAX - 6, -1, "beep boop", "e511aca4efb2443f75247f99086cad1c2e1b8aa107f2eec28b1e"},
	SumTestCase{SKEIN256_MAX - 5, -1, "beep boop", "0cb6fba5642cb25e42d098a05a96f6a999e3c36593b45c469dab40"},
	SumTestCase{SKEIN256_MAX - 4, -1, "beep boop", "8459b356c80dc769be0c599ccd461a1d23fd16bb6d4fede97dd32652"},
	SumTestCase{SKEIN256_MAX - 3, -1, "beep boop", "652010e976f2a479f976104ee840751469cfd6eacbd70892af314e207f"},
	SumTestCase{SKEIN256_MAX - 2, -1, "beep boop", "f8248d738baacb63aed745ea6a6a754106e1daa210255f4d5f8fb8a68fbf"},
	SumTestCase{SKEIN256_MAX - 1, -1, "beep boop", "9af9ec5914772e6b12144126dac2a5fe35dbd523a641394d0cbe040bcfeec7"},
	SumTestCase{SKEIN256_MAX, -1, "beep boop", "9f9156e984df2419deda0b70ba0a1140091b1bad7631bbc8e23d32aa0223debe"},
	SumTestCase{SKEIN512_MAX - 63, -1, "beep boop", "00"},
	SumTestCase{SKEIN512_MAX - 62, -1, "beep boop", "1947"},
	SumTestCase{SKEIN512_MAX - 61, -1, "beep boop", "48fd7e"},
	SumTestCase{SKEIN512_MAX - 60, -1, "beep boop", "2a1de6a4"},
	SumTestCase{SKEIN512_MAX - 59, -1, "beep boop", "1be5045f2c"},
	SumTestCase{SKEIN512_MAX - 58, -1, "beep boop", "1db6742d00da"},
	SumTestCase{SKEIN512_MAX - 57, -1, "beep boop", "ed3e78bba2316a"},
	SumTestCase{SKEIN512_MAX - 56, -1, "beep boop", "45f7b27cbfc594c4"},
	SumTestCase{SKEIN512_MAX - 55, -1, "beep boop", "d788d00508024caaae"},
	SumTestCase{SKEIN512_MAX - 54, -1, "beep boop", "d3d7e6d7ee36d8553ee0"},
	SumTestCase{SKEIN512_MAX - 53, -1, "beep boop", "e810163471cf6018be709a"},
	SumTestCase{SKEIN512_MAX - 52, -1, "beep boop", "44bd0e2c8fdead8774b43b03"},
	SumTestCase{SKEIN512_MAX - 51, -1, "beep boop", "c4aec9e1cdda14c52a2b913548"},
	SumTestCase{SKEIN512_MAX - 50, -1, "beep boop", "5ab1562fd543b39e5e96b428454b"},
	SumTestCase{SKEIN512_MAX - 49, -1, "beep boop", "99d81ef1ce64ed854c1cacd2c00741"},
	SumTestCase{SKEIN512_MAX - 48, -1, "beep boop", "1cd88c059a2883aca86c01dc0a030d25"},
	SumTestCase{SKEIN512_MAX - 47, -1, "beep boop", "d7bf31760b9cc001955d3e82efdf8708c2"},
	SumTestCase{SKEIN512_MAX - 46, -1, "beep boop", "fe1bbb7f92343efed4e4a9484bc599dc2b33"},
	SumTestCase{SKEIN512_MAX - 45, -1, "beep boop", "c4c0950bc20e73c9a874eaac44bf7d2a20cadb"},
	SumTestCase{SKEIN512_MAX - 44, -1, "beep boop", "d2962db2feb936a91bf2d0099319e8800da28433"},
	SumTestCase{SKEIN512_MAX - 43, -1, "beep boop", "58e834932a537fa5795723d4552271c5151dba56d7"},
	SumTestCase{SKEIN512_MAX - 42, -1, "beep boop", "808d3d226f5deb237beccbddab21e4f061dafc972ffe"},
	SumTestCase{SKEIN512_MAX - 41, -1, "beep boop", "a4f19769ee0a9147255b32be95e52677228b49dab8b352"},
	SumTestCase{SKEIN512_MAX - 40, -1, "beep boop", "6164b8f4f9978b3625df4dd0ad8e2de4e8b57f91c0e93f76"},
	SumTestCase{SKEIN512_MAX - 39, -1, "beep boop", "7c2e3c424780df31b5647018ce26222c74931de163d483d914"},
	SumTestCase{SKEIN512_MAX - 38, -1, "beep boop", "feaa1cdbd4d6eb201840cde63ffc58c80367f1ee072c0acf96bf"},
	SumTestCase{SKEIN512_MAX - 37, -1, "beep boop", "539dd952dbaeea70c9afbe20cf21086df442493a8dc26e8d2c42e1"},
	SumTestCase{SKEIN512_MAX - 36, -1, "beep boop", "3311f2e808898c97471ace98d3cf2ccee6b7fd8e2b0945f653a8bd37"},
	SumTestCase{SKEIN512_MAX - 35, -1, "beep boop", "f076e04ecf9517c6c41ed5621dd0ac7c14d64c1ad7c3be45d000efb975"},
	SumTestCase{SKEIN512_MAX - 34, -1, "beep boop", "fd22ae743aa09d2a678168e838110d1029597c7621e85e85240e597020c4"},
	SumTestCase{SKEIN512_MAX - 33, -1, "beep boop", "19b5a1233faf47b358237e2eceb1260028846a4866e3947c80e35e0db26b22"},
	SumTestCase{SKEIN512_MAX - 32, -1, "beep boop", "76964cd67694f35bdea508d88680acaba088424823a05834e280354b735f1a0a"},
	SumTestCase{SKEIN512_MAX - 31, -1, "beep boop", "7e6e19375221e77fb29505faf2d7389c1c489bca56aca8f9aa0a3fc1b1694f6581"},
	SumTestCase{SKEIN512_MAX - 30, -1, "beep boop", "5e3e2b63b2795f96a0eb8f5cf8ad2cdd6692cee376fa46e71786bd6b5bd900aca5ae"},
	SumTestCase{SKEIN512_MAX - 29, -1, "beep boop", "fae0ca21ae226b95967e045d22b940ece1b5a43ec63c6fefb5b39553dc4a31a6783071"},
	SumTestCase{SKEIN512_MAX - 28, -1, "beep boop", "a3b0e169cb0580f8044bb619bf98c666861281d966417d774489d5c78353d2ab5bfa764e"},
	SumTestCase{SKEIN512_MAX - 27, -1, "beep boop", "e092b455c23cfaafd1ac0ad14b9f9eaedee911fabb6713f131f75763990f92619546d5d2b6"},
	SumTestCase{SKEIN512_MAX - 26, -1, "beep boop", "7b091e347c196e348851851c469e3eadbcd2bfa778149e02ed68d9f939c9361091e39b9d7c6b"},
	SumTestCase{SKEIN512_MAX - 25, -1, "beep boop", "24220db3952eda1122e19c1c1a73465a61ab4bda6b6b0f563af4645638f32783e08ef37d9c7e01"},
	SumTestCase{SKEIN512_MAX - 24, -1, "beep boop", "c355b8281212fa2a83aaa32b52aa3fc3c22e1acdb59c815060a826f9c6fcbea9abf8ac5b9fdafaa0"},
	SumTestCase{SKEIN512_MAX - 23, -1, "beep boop", "481f5a96525d9b7b42cc6c38b76e735f46f82961b9f931683ab2474d7512eea24eb0eb32bdca485732"},
	SumTestCase{SKEIN512_MAX - 22, -1, "beep boop", "40bad4699ad90379966b3f99f76e3e22d850764571c4f654c43f38cc55e59adbd952bacc0de14e9b38d1"},
	SumTestCase{SKEIN512_MAX - 21, -1, "beep boop", "2eb2253d2234d42a615a20b6c4896f655eac0c891fee877b5fea79d1c230cd89020946d799ef2290e1465f"},
	SumTestCase{SKEIN512_MAX - 20, -1, "beep boop", "e85b094d6c297b5e8e7724b34df731b183d5c2addeb490995287005b6acc2fbe90683ff934567150901fdf03"},
	SumTestCase{SKEIN512_MAX - 19, -1, "beep boop", "da42a774f8a1701fa83d36e1debdd69b7397fd91630619ac80cbf44639088ad30a2100f152426e71a3ad5a961c"},
	SumTestCase{SKEIN512_MAX - 18, -1, "beep boop", "800ae02dfac79d2d6fd14e74601479602527978893f096353d977ee2bbcbd91a17cad13540b70ce88fab7c7d421d"},
	SumTestCase{SKEIN512_MAX - 17, -1, "beep boop", "0339759257e6f09499cd624aa8fc1819f9d0b210432d1de58689ffdd05e328c6e92e6e3114154f36a04581014642da"},
	SumTestCase{SKEIN512_MAX - 16, -1, "beep boop", "80c7c4940d29c3f56c63776b96f05bbc0966afb3100cb8c573a5a0766a95112c6d550ff0d106f849de934fa640e25b95"},
	SumTestCase{SKEIN512_MAX - 15, -1, "beep boop", "a6e6bf2d8b4c570d57985620be85c2b6a842592fcccae860e3186f3de99cab0d2f282aa002bb4beb793a59554cc5e1b492"},
	SumTestCase{SKEIN512_MAX - 14, -1, "beep boop", "baf605f41b9ca03a5515ee35ab20f9e102de12e50332c5e43f909c7dded9c7a7e497c23d2adecf71e8f988277e032001392f"},
	SumTestCase{SKEIN512_MAX - 13, -1, "beep boop", "94cd4c5aac9f48778b469ce3b8f9aea81723f13c9f91d4cdde719091f98f0d41046dde9e6b4be75d6f60ddfe73313afe33d793"},
	SumTestCase{SKEIN512_MAX - 12, -1, "beep boop", "4d599ecf30c4df3feb9d06388b79ee1938e5433a9a0b6b8b9f6d65023105fc1f7f2807e997d7967102c33b2530f59a0fb4cb7a09"},
	SumTestCase{SKEIN512_MAX - 11, -1, "beep boop", "5f76b9b7fe05b20f604d203ae893844b0d3e3d9dabd18c1085780f147eaaf116eefec7c277af3d3ec202fa03caf9f09817a6ffb130"},
	SumTestCase{SKEIN512_MAX - 10, -1, "beep boop", "dd043d995ee17ff1e80affe727eaae999196987e26d3a8afc480cc2a4d76fc32fcfa09e63f0e73445eabacb380d7638d805af0f43900"},
	SumTestCase{SKEIN512_MAX - 9, -1, "beep boop", "c3e0cf4d422b2bbc22769859d8fa095d3d991917be7d446c4eda7b2aa637b0f7a1699b347d51d541153d2b9f4553dc4e16f19e5affa046"},
	SumTestCase{SKEIN512_MAX - 8, -1, "beep boop", "13d1af5e5e35b1659f22cdb9a58fff5c345e86fcc4610cea11d790133edb4aed684a3e9a77f7f9c93dd438c55c5ea9f33f9fb95b029a8322"},
	SumTestCase{SKEIN512_MAX - 7, -1, "beep boop", "88a6d7fe89ce1962bc4e93ca1e6e690a64a881cfb9726ec998570d730dd04a439a7dc32378fb92a18761a948804c61fcf185787fb96031edff"},
	SumTestCase{SKEIN512_MAX - 6, -1, "beep boop", "3ec1443625e32419f7d781813c9a9ab9d0e43eddaebc6437728f2c9fb8633b6804934c66a3832480facfb92a055d49f72e63c5abd41c11f527ef"},
	SumTestCase{SKEIN512_MAX - 5, -1, "beep boop", "6d97ed1818ffeedfbb189cd3c16902973e8c2214f73a0027e1e8430c80ce31c2bd2bdc33c9b2175c8c06285ae0aa4c20aa60edd5f86d489917ee63"},
	SumTestCase{SKEIN512_MAX - 4, -1, "beep boop", "2bde247fa0b9f200ac289ccdd4e206f32f04a1b9212720de1b6daa031e824bb1807d740e6bceb923a1aa263e5768ed347ce3c4da5d30391ac45c9d84"},
	SumTestCase{SKEIN512_MAX - 3, -1, "beep boop", "586ce1590fce74c60fdd85c068675354834d72a6cda221780ebacd57f4a1e9c4d6b3cb3528609a9458fd6da72f6e39fd34b877b5c7fb3c4942e0053716"},
	SumTestCase{SKEIN512_MAX - 2, -1, "beep boop", "ca755cef50908d5b6f62b28f5077dfd2dd7bd12620bcd9930059e87535c22f5d0c8325175cab780040632d44bc4e6f1a663c1f13f42e7de30d65c06753e2"},
	SumTestCase{SKEIN512_MAX - 1, -1, "beep boop", "9fd9b68694613d7af71249cc40ae6b74b4f37f05f9296faaa117601ed9029d5fa29406e47452446dceb487c04c4c78b1478d77bd583b74bb65b0c3a83ac06a"},
	SumTestCase{SKEIN512_MAX, -1, "beep boop", "949fb826dd57518665c641767e3c50b9cf1279eeea765bbc6f6bcc5f4023df54ef100b492142737b80c74843ed401a475f1923e2fc35d5ce2fb4c6daee5d1d56"},
	SumTestCase{SKEIN1024_MAX - 127, -1, "beep boop", "f5"},
	SumTestCase{SKEIN1024_MAX - 126, -1, "beep boop", "0cbe"},
	SumTestCase{SKEIN1024_MAX - 125, -1, "beep boop", "ba6964"},
	SumTestCase{SKEIN1024_MAX - 124, -1, "beep boop", "faf752f6"},
	SumTestCase{SKEIN1024_MAX - 123, -1, "beep boop", "2e3482aa15"},
	SumTestCase{SKEIN1024_MAX - 122, -1, "beep boop", "61e93ec88f34"},
	SumTestCase{SKEIN1024_MAX - 121, -1, "beep boop", "0c0b837fbe32ba"},
	SumTestCase{SKEIN1024_MAX - 120, -1, "beep boop", "a55454bfd6f808d2"},
	SumTestCase{SKEIN1024_MAX - 119, -1, "beep boop", "81235141d3565ec8f5"},
	SumTestCase{SKEIN1024_MAX - 118, -1, "beep boop", "666410264eec60246709"},
	SumTestCase{SKEIN1024_MAX - 117, -1, "beep boop", "307448a4b8ead26fa49601"},
	SumTestCase{SKEIN1024_MAX - 116, -1, "beep boop", "19f785e6670257e05ebf7be1"},
	SumTestCase{SKEIN1024_MAX - 115, -1, "beep boop", "ca7dc55a211374ff605e164641"},
	SumTestCase{SKEIN1024_MAX - 114, -1, "beep boop", "4393b8a142454a878c25ebd521ed"},
	SumTestCase{SKEIN1024_MAX - 113, -1, "beep boop", "c683d762aa781cb0a29f4abf88a688"},
	SumTestCase{SKEIN1024_MAX - 112, -1, "beep boop", "4ff0939065a972c1ba74e2b9382e5665"},
	SumTestCase{SKEIN1024_MAX - 111, -1, "beep boop", "6896f7a3f70c4403da3ee30ea6765e7b85"},
	SumTestCase{SKEIN1024_MAX - 110, -1, "beep boop", "e4f38839da1bd3ce949ba4dbce8d1a347587"},
	SumTestCase{SKEIN1024_MAX - 109, -1, "beep boop", "12c050b013983d432a8df4017a199c53719474"},
	SumTestCase{SKEIN1024_MAX - 108, -1, "beep boop", "ff26364fa85e2cbf7d1a21af746a13567f2847f8"},
	SumTestCase{SKEIN1024_MAX - 107, -1, "beep boop", "632f7bba8d0e1d9d66309ae974dc941adb70962b56"},
	SumTestCase{SKEIN1024_MAX - 106, -1, "beep boop", "7241ad736eb8ed718919e554e546ff2a3413095af65a"},
	SumTestCase{SKEIN1024_MAX - 105, -1, "beep boop", "5702f425035067e49ac16861769087af12b43c0f5a4315"},
	SumTestCase{SKEIN1024_MAX - 104, -1, "beep boop", "117de1ffb97d2c74f328091ad366ba2f521a4dd541fbcba2"},
	SumTestCase{SKEIN1024_MAX - 103, -1, "beep boop", "ff0c09478bf999a431d96a8a4d84da633413a80ba94770a99f"},
	SumTestCase{SKEIN1024_MAX - 102, -1, "beep boop", "a891b28b5117c76e0b9aff9def3bdb06470d333be5a3114e3274"},
	SumTestCase{SKEIN1024_MAX - 101, -1, "beep boop", "6b0fbf409ab601e3738f4b71cbfa91ff18ac3d27f9f46b767a8129"},
	SumTestCase{SKEIN1024_MAX - 100, -1, "beep boop", "9f49d64cb4a77908b22ae7b1be823793eef52266a381c2231641014e"},
	SumTestCase{SKEIN1024_MAX - 99, -1, "beep boop", "adb816e6176f07d7dac9cd530219be6ef650026a32d3827712906265d3"},
	SumTestCase{SKEIN1024_MAX - 98, -1, "beep boop", "2b7e6b7c5fab3184452c1866b07b0316192cc885031a50f41759b83c0901"},
	SumTestCase{SKEIN1024_MAX - 97, -1, "beep boop", "1e64032ed0084549f93562e647a61f675be6da3c80c24a5f07031b07b37d2b"},
	SumTestCase{SKEIN1024_MAX - 96, -1, "beep boop", "c8dea5c497ea660e20c2e3b7e9f5401fa3c715fd63d437fbd94330cc461d4e2f"},
	SumTestCase{SKEIN1024_MAX - 95, -1, "beep boop", "a076c68e96ba91d950a82fdbbcd23995713f58c99b796bbc8d0f2121f4d15ce9e7"},
	SumTestCase{SKEIN1024_MAX - 94, -1, "beep boop", "d70ff21ae5b4afe81961328e08e0c84ac9ee0079154aed4cddcc106cbf777c0c8bed"},
	SumTestCase{SKEIN1024_MAX - 93, -1, "beep boop", "38cd69478be4d76b7c317a903ab279170d0835e6aa45ea6d2d5495cf2b471fd8fe4dbd"},
	SumTestCase{SKEIN1024_MAX - 92, -1, "beep boop", "63215cc8133dbfa42e01dacae605666ce462248d9a7ce8cc55c80e5ae8169a5d6d9c6114"},
	SumTestCase{SKEIN1024_MAX - 91, -1, "beep boop", "96acb5cb75564cdf031c59d4b46bec1f88f003a792fd3c20a859041035f883dbd7a95644cd"},
	SumTestCase{SKEIN1024_MAX - 90, -1, "beep boop", "d423c6dbfe24d474c5a6b3caf81096b1c57399adc94f71ea40e5aaed456f847b4f4616bda266"},
	SumTestCase{SKEIN1024_MAX - 89, -1, "beep boop", "b6a474d3edd4e6df275847d4039ca587c58750c3fbf84f3c29d4fe8cfe1951b812af01636aee3e"},
	SumTestCase{SKEIN1024_MAX - 88, -1, "beep boop", "dd4e2d1bb09f7b695f204f97a4809abcb69327ec8be3bffbb61391d5f6b96c5268993f83b452b1cd"},
	SumTestCase{SKEIN1024_MAX - 87, -1, "beep boop", "3aab8e05c9f3bf214c07af1fce42d8a787b9e853c23bac9ba4e667c0b025d406a69d9ec70a81d8d29c"},
	SumTestCase{SKEIN1024_MAX - 86, -1, "beep boop", "3198b1e4a31ed1e2823f8d6eb32aa6a50acdf5e4080873e7913d89fc990fa32abfc9c41ae55705ba5a34"},
	SumTestCase{SKEIN1024_MAX - 85, -1, "beep boop", "35491a9be939401d12c395b2527431f72d7b03c1ac6c46860a833c70b7ce5b24996c92221d223077a1b17e"},
	SumTestCase{SKEIN1024_MAX - 84, -1, "beep boop", "498e52323f6f8787f8393d3135c6d937978d575f5e79aabf3d2e52bec7e89300be019727161232df3ddd02a0"},
	SumTestCase{SKEIN1024_MAX - 83, -1, "beep boop", "e54cede397510865cb08ba2d5e03b42f4281e4fcc3e9bba78e4b5a649347bea24122f2d5b2a0c31d359e9d2925"},
	SumTestCase{SKEIN1024_MAX - 82, -1, "beep boop", "ad2c77c54ec919a5383d5b6b9902d626b695c8e3aa521118a803dfde0db765693d150eec3ec905f7a9ff1b790ff2"},
	SumTestCase{SKEIN1024_MAX - 81, -1, "beep boop", "2026e035e78fb7a682b99497caa254d28693331564bd4599e1e7cb9ade82e6e477f94f1f62b7bfe458799a6dd3e7e2"},
	SumTestCase{SKEIN1024_MAX - 80, -1, "beep boop", "f931061f58d36fcef19e69d3ff63095fa038eac003634e5504283ce1efe18d5f3234a06ecbc36a1a83f20f45f887888a"},
	SumTestCase{SKEIN1024_MAX - 79, -1, "beep boop", "9be870cb114bbba3ea88a81c55e03f64ab187f30d952ba3d4ce8287809249b85138d44c388fc615637727e5aedafbaa43c"},
	SumTestCase{SKEIN1024_MAX - 78, -1, "beep boop", "4d4c0b0275ce783752a8d8d50f7f73ed519445260527a2d5a7ca530cf61c645517108b1c9263873a1671400af11ddcfa1cad"},
	SumTestCase{SKEIN1024_MAX - 77, -1, "beep boop", "460ae4259725e31cfd6caa1d60ef88168a4e7251a2e9efbf89bffb565a54634c2f71419ad88ec66bcbfbe0a7186a7544dc5656"},
	SumTestCase{SKEIN1024_MAX - 76, -1, "beep boop", "a865f91f025f04191f9956f7c24600cec23a0848d716080bf9210bd7f9075246657390f35356d5bfd295866d6908a572c00e7a97"},
	SumTestCase{SKEIN1024_MAX - 75, -1, "beep boop", "7c9aab7d230b9bad86d39f73bd164e50fd8aeae5e2687176730301b0400822de7c918096e088ff7d36140619c2ec88a9def4c1e486"},
	SumTestCase{SKEIN1024_MAX - 74, -1, "beep boop", "02f57498fb92fc5566544c17cb0e05f181e1fe3049e8e93ef6e32a8fa148045cdfa9e83fd1dfd885ab1022f9f5a566e5a0963685a948"},
	SumTestCase{SKEIN1024_MAX - 73, -1, "beep boop", "fd6b6b776169c0a7ac3ee82e4271e40d7b8c4e477eac79ced4781c1c2d1f19e2db059601b480a5298ca88c0794d228a9cb5d66e24006c0"},
	SumTestCase{SKEIN1024_MAX - 72, -1, "beep boop", "8b33749fb7b735f0a449670da8db05081a523a4b8d7fac391f2a87d8e52d3fdfe528013ff17ea178171e6180972ac5cb55939baa213f19ca"},
	SumTestCase{SKEIN1024_MAX - 71, -1, "beep boop", "fadcfd830d44c111db222442a142bd384b7628f1940bc015bd2016aa74584c6202b6694c29f81d4dd5f5a13d64d34d0d55e65544814c1158ca"},
	SumTestCase{SKEIN1024_MAX - 70, -1, "beep boop", "517979037fb5f49a6d3e1162ff4b46e5a20c98236c19c587faca21902e8e704b3a99a66c27b6d370d05dce5ebc2708a734a046e9e2389c879f00"},
	SumTestCase{SKEIN1024_MAX - 69, -1, "beep boop", "b50b43426daf7dc1c93018d53728cb8580650f926da8c40bc49700ec9587d5309a23e0afd845d91e1d26ccfe3b38b44a0c30d96215eb112c1a6611"},
	SumTestCase{SKEIN1024_MAX - 68, -1, "beep boop", "dcfa838b4541898895adffb6854ff1f526b85a5c925e858357e747578e811966e32e8f17259a89e9ff90c71136ec422c1ce2f24def9950a645f5c290"},
	SumTestCase{SKEIN1024_MAX - 67, -1, "beep boop", "d5f061c455f3b586132e7aeeb375739150839bd2cc980614d7fd187549af4eac867b04de0b407f15271db77df7eea9d83239082093add00fc02e02aa2a"},
	SumTestCase{SKEIN1024_MAX - 66, -1, "beep boop", "0c6b19cbf0b1a1292484449be7c5b61710ccaaceb0a78b0ba1a3aa0bfb437e36c804d5f7ce018010a8c752edabff51b1b7fb81c6a0a82efdde1af1b54c27"},
	SumTestCase{SKEIN1024_MAX - 65, -1, "beep boop", "4fcfdd0816d9dc450e5344b65780d463536793d5d266597ea571b54a445216f7c6706def3c95fe7e9907e61789e5e9a9208d6e2d64a1376cff65e33fde454d"},
	SumTestCase{SKEIN1024_MAX - 64, -1, "beep boop", "abb528f3880534aaf29016866d01b94db7d85c7529ccddc51a6beb7a670a8af9e3635e3ffac0be4f99a644478798b25a6e0ba845d34c3417c62dc7d97d6a58c4"},
	SumTestCase{SKEIN1024_MAX - 63, -1, "beep boop", "1acf162a2f4357449261720275d519190dad6779b1c680a7a32bfd51e803b7536162681d3cb3646816e3c9cbc7ffbb9afd7b90c0ddf7f0db243f7047f713aaf3c2"},
	SumTestCase{SKEIN1024_MAX - 62, -1, "beep boop", "3f16171ce952bf1dbeb8fc8554be3aaff37c964b1ea7a07f37e74770ca82ff9234a263692d695bf6bc391c817433a856cdc9b49af5108f1c209a074d284b27d65907"},
	SumTestCase{SKEIN1024_MAX - 61, -1, "beep boop", "3ed8ca7f8bd585e7ff475e238c0084db57aab90a137c78f80b8ff4eba543c9668b5c4cdec012e1c89d4bcec54286bdbe81e3a4dbb415972c9aeb6aecbe6b139e8e385f"},
	SumTestCase{SKEIN1024_MAX - 60, -1, "beep boop", "9f3448ad1cd3c0de6e5082d2950ea665439dc754839ff0156cf1e12831efef108f02b7801edae79f87f86488f9b95c8af215a8c15960bfef988fb2c6f3f02399aafcc6ff"},
	SumTestCase{SKEIN1024_MAX - 59, -1, "beep boop", "e785267715c65d5d7ead7049187c1e308f0abf4e63f8e708203a7f6c5de713284711589e88f6f549e2836f37eeae7fbce899347bb625309a1a07327dd055a5c53e7bf9d3f8"},
	SumTestCase{SKEIN1024_MAX - 58, -1, "beep boop", "081aa6a7090a1efaca677a6372bdbcc48fb4deafa552e5d1f4526fa4a1615859d0ca8ae77f45cf066f72b780d7b840913f9e27e5963cf8143c1bcdffe4b04da20ab0368c0053"},
	SumTestCase{SKEIN1024_MAX - 57, -1, "beep boop", "4af768a2eb34916c1ae63473d9862c9736449cf81975ef81385b6c0a9a405755f5b042ad52762b4221e921ba3e6bc9e84031ccc525fcd4567bfa03088f5e650990f8802d4a76d9"},
	SumTestCase{SKEIN1024_MAX - 56, -1, "beep boop", "980eea698e3521f2fee8bd354dcd50d3d322c6979e13a71ab505cd26b33f77bc6456853601174a081c55b6a5edce27ba64056ba21f606fdb79e6a3d78d952af8d0482385e16860c6"},
	SumTestCase{SKEIN1024_MAX - 55, -1, "beep boop", "3c1ca271d8e1034983211b29255805fe0eaf3e68fd7a1b1acbecead4b58a69d02833f9f964840607a37abcaf2b061c1407f0bb28a568b03c25b47e3ad57dd52001847766e5355eab25"},
	SumTestCase{SKEIN1024_MAX - 54, -1, "beep boop", "9a5b9ee3e2309b195335347303b07c81ff502c9840f5476cf41aebae5df552981264a688c25d8d508b81b314dc87e2213172360be6e6b01d694f720755d58bbfa3a9509112aba6cea936"},
	SumTestCase{SKEIN1024_MAX - 53, -1, "beep boop", "d65e373f8a142120f2c8c38f4d95532ce2a19decb28f97c94fcaad807553a697eac0ec833c7ffa633268da3ecef8f2ba450a30be150a769d61c56d36deacdf6234d11deeb72df1bf06c5ff"},
	SumTestCase{SKEIN1024_MAX - 52, -1, "beep boop", "9222adfc61eeb45785fb29b39018d83d4c1f3554016397fce22bfd5933884785c1edc18bd2b08bd414530d56b32c84b86eb399d45ba99d40070c5c011c099b68f90fdd2fa37fff4f29797e32"},
	SumTestCase{SKEIN1024_MAX - 51, -1, "beep boop", "b11a8b4d61b3524cee3ee63954a3270837678ab9a67720a7ba087a0f9b180f30595fc881118b9d1bbc3e8d210af76ce9c8bb7ac39b6d7b8c18e20ce6302c634db41e74c82255d0b7e062f45a00"},
	SumTestCase{SKEIN1024_MAX - 50, -1, "beep boop", "104eb4f850049c6ca93ea34e648c466060998460861a6e6eaec2a82ace37672c7734d1302cfd5bd7195e016643a4a98e3611e600d4cbfe0a06d606f3d863e8a7ffd5fce6c5ec75d163b8b106dfc9"},
	SumTestCase{SKEIN1024_MAX - 49, -1, "beep boop", "a9588b854918cfefea75723adff6ef49fe5a690faf2e7515a7f80c4f3d02424cd7b591c39de984a04a87c8b12b695c289d887ede72e753af4f59956c7ccd0e600bd086aad9fdedb7270fc074bb5cf0"},
	SumTestCase{SKEIN1024_MAX - 48, -1, "beep boop", "ecacaa021b5f2095ad008d544ed8a8786aa5b1ef470fac3e8fa0848bd0edbb21157912d1ae0691a71ef4824f0b50dc7999ebe834e6c7670672866ad14851d6b9aa6c81f3e8a2f668f35bf17dc31dff3e"},
	SumTestCase{SKEIN1024_MAX - 47, -1, "beep boop", "af36634df47fcb7eaafc6992978ee006cf0797df6a350c4cfd47f05973d8e6534c0d1544acb9efa8e520ef3220efde074b01003d1695ea8cf586b02180bdc3f0b9237e490d15d0452d91a9d39659b51984"},
	SumTestCase{SKEIN1024_MAX - 46, -1, "beep boop", "6a2ead45cee5c2282ae3f33e2d2d3d34d6d11f0e57e3c551807a486b348ebdd844843add516fe35f822c543142a5a85797bc6eceef0f1f766c9fe668902be6e1fbe8f75bb8e36413fccc845f99493861cf35"},
	SumTestCase{SKEIN1024_MAX - 45, -1, "beep boop", "0965f2dcc471633333ce83ed619671d73f714782fe81b5e91d5cc4f38e456bc28fe3ec883239b8883c358d567ea3e8950f20621c86e7111ff9f1d0f446cd7c73ffd542172759bc760e6922e2a99e557b4bb6d7"},
	SumTestCase{SKEIN1024_MAX - 44, -1, "beep boop", "3d6b9e34f54e16b5d576b55e65207b55234fd5ab3993b306abf66e005d61976700f6c915dc5e056ea650cb2e470416fb32ccfde8a0aeb4baa4693c1f3eaddfc292136f4c216d5ee2fd7f73e53ff501cac45e0fb7"},
	SumTestCase{SKEIN1024_MAX - 43, -1, "beep boop", "84d8387d3b1cfc904d8d1ca5ab903735925653046a8333315dc2f86c64569f8108cc1690f5f4f608cdde69ea9235a8cc0c90679360ef94c02b4c578c228b3aaa4708f575037aadde235d110964f3826b31046bd5f6"},
	SumTestCase{SKEIN1024_MAX - 42, -1, "beep boop", "758d972830e014dd97d2b00788c0bdc47b294f1ff748e940625e41c76895f4b7d0346620d132dec1228c70c68ab43f655f2017782c9571cd7cdca5c554b9975112a9a875aa977ad030e0cf4c2021cc8d6ea39eede995"},
	SumTestCase{SKEIN1024_MAX - 41, -1, "beep boop", "e30b68c49f155b168af31030023d0baee810ed3552e8601b5473a0003e0ed9816f4010775d09cee6fc2cb0a24faf869a5d11a31fd4a19261494edfb394c114bdc70688cb05aeb384b1b6712cf0fd2c9bf7b6fc61a3106e"},
	SumTestCase{SKEIN1024_MAX - 40, -1, "beep boop", "753e432acfba2d905dbe69b57f29a4936519aff342e2ba953fcd1fadb48f54fd73d209d5a89bc534369194c4b1a3bbdf20763e414a0e226b006b465d1a32cf2d29c3a54f9527159dc2d06fecb14216c557a5a74f78f8840e"},
	SumTestCase{SKEIN1024_MAX - 39, -1, "beep boop", "d46f2dba1ad54a86d13884357203342c70342e4133ad177edb3e2121cd8d712a4edab687d77fdd2952162d840a5e043f343dbe1e727787356ad03a09c589c078dd88c0a63e6f45bd536e8bf4c3120016d6c0b3e7820131b2d7"},
	SumTestCase{SKEIN1024_MAX - 38, -1, "beep boop", "ee7278d2d90ee33c443cda2bfa3bf031e64d6cf01331c4e712a2aed258e53332f29e504ab4ebeb06309ae599581738d8e68ec7ca445c85f7eac5ced94b2f06be5996a1ebd59330704b534d26dfc4ff759cb061ebbd86dc7407a8"},
	SumTestCase{SKEIN1024_MAX - 37, -1, "beep boop", "bacdc5f96ac19992f16e1f8435583e14f46c716ec07bc19f51beea64f8982d9a638e030b4024ea6e4988c69e5a656b172022b425a49823e7c597e4d6227797481311cd3b6138d391d379e3dab308d13a3e2cad3804927967a9dbd4"},
	SumTestCase{SKEIN1024_MAX - 36, -1, "beep boop", "10b5b97a03637fec5fc4f46cca784603a894a22ce5ee9cafc3106844b437929b2d365dbabc005541afa07861b2073fd154bd3d414a12f3516be3b7d68250290ad5118b7e3fc6312f17dfb5b3fcfc0f3b529b0b641791c5264c50d57c"},
	SumTestCase{SKEIN1024_MAX - 35, -1, "beep boop", "fe56c12e76e42e50edcb876b701ca84ac203a406960faf6ab87564648c4d8fb55473cdac96cdc9fde9f228663822c8cd397a0f417a4b5b66dfe9e5765fc3503f8109d305cfe4d8f1d09ad8e2a7ce23e030d7112addeb01277447afea92"},
	SumTestCase{SKEIN1024_MAX - 34, -1, "beep boop", "11052ec1f549fa2e0d88faa4fb69f05da6b73814497d4a253ce418ad8cd24f36446c2385df631fed60e7651a6e4803bb61cb2e051af19dd68c2309c91c3b1bfff75dcfa3e3be11435d576cdb37a2ad3fe7268714971ee68f5d1e14a6eaa4"},
	SumTestCase{SKEIN1024_MAX - 33, -1, "beep boop", "4510a25ca8c4a6ac17a924915db7dd70f41cef14fd74668bd23253b2b88adfe6be1e3c9a26225a6038a14a018f1d3c109c3ce1f7a0e8c041424f0e4875c1de1496ae8e93f6063297b460605cd0ecc8f1195ab3d07d1473b590710635b0fc60"},
	SumTestCase{SKEIN1024_MAX - 32, -1, "beep boop", "66e3c616e72ee146b2b94febb28458585f8efc8b53032106db04441fb078c022ed22dadb981bd72c9f10cee1d081285f53fbbe8f76ab53d34c5e83c686ac799f4bbffa1d722ee1e2e7d0aa8e8756f92a69386ff6732d44af007ef9167c82aa5b"},
	SumTestCase{SKEIN1024_MAX - 31, -1, "beep boop", "959db77d9f357ad42c6f9999e16044ac28fc303842207caa595244e1fa5c8fbaf5988a219699b67215763b5dc6db4f681cbe8432c37664a7c69e55bf545adb329b984159cbc2e6f161422ff52779be02557375b232d202c8a771d2c9d4886220b3"},
	SumTestCase{SKEIN1024_MAX - 30, -1, "beep boop", "7d20a1a13b30c281d45ab9b251dbd04b1d0c10ecd5fc7ed4150dc9a15f6798036c83be48814b04d6dc8f2f9bf7347a6eaceb4a70302ae51f6f716cd3f64dcf5fcba67a7b8c788fd9482a9eb245b10d602543766a57533f743e939ece53025d3ed7bd"},
	SumTestCase{SKEIN1024_MAX - 29, -1, "beep boop", "7ffa2a2d01e477f525c18f7747fafb43a669385bc487463d2f59283dcd3052d1fb28fde4e48f278f926b82e5ef6308200050e5e093fcd8ba5a14a75681a3ab721f348c375630d93c31925a573991aa95ef4d55990110a479fd582a5ae52b7e85fa6877"},
	SumTestCase{SKEIN1024_MAX - 28, -1, "beep boop", "b097dca1930b35ae170e6ef3d21ec356eb11a6e681879ba34d9d13e2f7dbf0bbccc8f77f2a38e5cafa46461079b2a8377255ed4a5e047de2e6b0c2d3fc66f8573eea3663888214065eea5c3ba3869d199d36072b82a2d91afcc01c4f69a2ce750461d8e3"},
	SumTestCase{SKEIN1024_MAX - 27, -1, "beep boop", "f48321be48bc2f25c2ac0f08431d4c29a588ab24d548f707d493c52e934771f92db36da8cbc6933fe42b60c3316a8ddf782f4ba6bebab0ad4e0db256959866cb95be83c7ff5411475682f553c41ace0bb0edf81a5f2288dddd3bd23ed198946e9f05e06021"},
	SumTestCase{SKEIN1024_MAX - 26, -1, "beep boop", "7ff0c191f1b8caa11e1270f4f9547bd665b1ccfce3d5ee95d39965a450049f2d63df6fe0045ab63271b8641e407c7edda1d170389b3c893cfaf5a7804ddc8cc5b6bc37aad3df50a6a04bf5b482cbec53c6cea0dccbd2d7aaf7cebea305ca9377b1c44418f154"},
	SumTestCase{SKEIN1024_MAX - 25, -1, "beep boop", "d283590cdbf28b92ff1bbac4693da6011afbb902fd7591788f7696098d8dac63f6f59bad7e3e0686756ef1cd74ea955e9e1e9223a6ab12be49fbc64191d5ef59b54f3d75ed1dac56b15b9777b94b101fefe6392632a379e585380dfef26ec76be8c05a2f3462df"},
	SumTestCase{SKEIN1024_MAX - 24, -1, "beep boop", "186e30c80a8abec32707c63a007c8297ce3e369757da1ac168ca3a9f7f149edb44b26085b68d70ba888046bdd7df9f0aa655138e6ca9aee5fd473078df8a91cbbaaba3b57190f18dab8b8cacb352dafa881671ac55206029156c243dc8aef15b6da41105bbe92701"},
	SumTestCase{SKEIN1024_MAX - 23, -1, "beep boop", "82508c13512896abfcc99fed9a57ad53923f8a9667579c477a5ac1f87d3910b5679b59287fde639b47d97d780335affda43b4df1a0481e66c1709f2077130f84f65812557e4557c73e5d0456cbe08c2b2aa989fb6f535e41e6118a15ac491ad38a9d03b638dadd42ad"},
	SumTestCase{SKEIN1024_MAX - 22, -1, "beep boop", "76be019bc0940d56661a8ce57592f11ea00279876541a719db6ce768518da5bfb341c82acc05b8e85748cb98e39e6ca757f43bcd184c8ce8f36c5b3d74cd60bb3ff48552f79853f6085117418519d95a2905942799997abdd59bcbb7da8bad4b5b154928218e64020001"},
	SumTestCase{SKEIN1024_MAX - 21, -1, "beep boop", "1bbcd8a61a0c885eb24848abe569941daf59d354b66ef3b819ea9fae5a9c0e5a98102071f3037b998da6ce9f5371c5864b4c9f5289d46625157ac77265540f655891ef28a26748d49c8870a80667336faa5973a945a86eeafe202a92877b949cc9c7a137775b46cd2ffc57"},
	SumTestCase{SKEIN1024_MAX - 20, -1, "beep boop", "d50d321f8dac0dacf154ab62df3c38142f511184cc7c43f671069481abda3f617b3facbd6c2bf8cb0376a89c56d71d8ac8676e504761f232786c8b5f34eef0025f01922fcf52acdc7785aaedd937502726eeeaecb91703bcb53855d8614cf63ccf66d8dba9342dfdefd07bf0"},
	SumTestCase{SKEIN1024_MAX - 19, -1, "beep boop", "9343199fe3f5d7d5ef43b6ca2f46e1284e20c45fb2bfdb893fc66d39134e8a658b0f04f5257b84b74c81c2e7f7ea8eaa5bad182757bba24332236d0120b6a72a4db4368e1da7a36ae36f5d79ef6c943c582e4788f41a55184197f949dd5ef47ef0046fbaaa9ac46449c6596ecc"},
	SumTestCase{SKEIN1024_MAX - 18, -1, "beep boop", "1356fbcb11a6d0f0acfb58271c1aad7278f2e5a4785cb0dcc9b6e2846203df8d68dd72a5c2e7169cd16d79e51c1ea59fef580ef76199a9a3583219e0f099e8013fbe9aad279a240dfb0e86c0011fde631a6808919c6c4fc6ef01488e48ff57916baae22d14a61ad9907b7ee71934"},
	SumTestCase{SKEIN1024_MAX - 17, -1, "beep boop", "aaa7b32b975e1890b1db92c2025e9f443fdbe3272f33bd9070a0cc95ad85c66f427e50217031e3dbbadd7d8dec911d8fffde7aedf0987b5e82a19bbea4db4f7164a087c1f3ea7bd86c56877eeb34a25da7762186c6869115189732965bf4fb991ed46712120f9143e19b936b52cc59"},
	SumTestCase{SKEIN1024_MAX - 16, -1, "beep boop", "590c1815aeef48d163cecb4fff49bc25069ce3679e1a510636c77bcf5a84ebd2f4161dc49434e8d191889c593596ed3acba761a5d0ef6fd0620bd3e1d98a9973086e2f30e009d29311c87efe54b56ef8a729e43930efaf623a4329a96fc136421143f4701e821ff3dcd015373c6787a5"},
	SumTestCase{SKEIN1024_MAX - 15, -1, "beep boop", "961bdacf479d4764cdfdc5baaa47e448800355ed921d0d450d9e9d5fe20ca480779a8185b21695a7fca3af573881a57fa0ded8b04a8853e1c8f864bdc81e32dd2e48418ce7c65f08ed7f5f60bc636197b295632011cdb71309313d8203dd17567a9e4e5653899c2b61fe222351081b3179"},
	SumTestCase{SKEIN1024_MAX - 14, -1, "beep boop", "7cdb0310e90f3d7b2f19d4b58f0ba5a5426abc9b5c2cad8684ee65aaf19a7bdab46758a4b8f9a6ddf204656deda27a785180a58a1757fdbb58ae81db330eef133e0740964fabcca6bc0b4efeb9feb54e5a4000779d6d39e5ef8b030c65349781e39b7bb61a143eac2a8f82c40a5ab57a1cb0"},
	SumTestCase{SKEIN1024_MAX - 13, -1, "beep boop", "37f856fcd8933a962342865e3d4343a8534c05f3ab33ef59af8f97c8e509fe20a1a8d20e8ee0cd492aa8343e317d48cd29b6766d43ffe67b8b38e2b63fc57d18c4d899e30be294d86db46eb8d1e6fc014786e063ef1dc1076f0cb20afcc4dcf8137c5d689f85a558d633db7036e2eda39e1e76"},
	SumTestCase{SKEIN1024_MAX - 12, -1, "beep boop", "b29fe6ab34fa4887a0306b692df937c18722167172d5256a30772bd60ce8e1ffcfe7b20b2d1cd6eb6436e91cc6281bf25d23d49d5e941defc514698dafc1bfb75fad2c0fdef305c0435d26386392d8a07edba9e9af6f553a034f648efb6f6e620ef5ef01e59404523248793d931185ee3ea1c28d"},
	SumTestCase{SKEIN1024_MAX - 11, -1, "beep boop", "143dc4027160996a7568da757671de145109985524e425acab31dd9370ec13cdca7f662af9f22434ab83a772e64c6559ee2958b03ac0bb94ba9d03fb6b57e63150ec8b2762a358348f9c8e7fad75d686b52bb33f35ee6dc051386fc080aa84a1aa92d3982b4101548da396bb529798f5d0215cb681"},
	SumTestCase{SKEIN1024_MAX - 10, -1, "beep boop", "d945537727b974b6142e23fa6ba8d5fd3c9a914295da357800e3ccf5a4ab3d0e3d686ab491553c74880ebf6c689e6942f8c88f587d190736060ebf9984e3e260419b98986fd91e44907bfa883a47eefcac35c19cecca740d3e671b2fe85bef6fa41afddc92dfbb0fefe9e624b4831eb0edc2c3e9c839"},
	SumTestCase{SKEIN1024_MAX - 9, -1, "beep boop", "b7fdcf0c832d40e79cdf69a93fddbf39d68038d4ccbbe81de3ea305d5ce0898a60170fe55f1966e4035c6aaaa6c7cff3911021804e21de971e547157111e3e526136211eab686f2d0109ce3a7d01c7c54da38a0766856a2e5e2fbe2df6678b1c979fe40fead4f713de930380e81809fde9a3e1dbba5c0b"},
	SumTestCase{SKEIN1024_MAX - 8, -1, "beep boop", "6b321bb810712a22035e0b1c34f69ec21f88529eafdbb3596682722af7373d2c6e832484667c9c7bb23c713488e861c891f982b300771e822c53e52d09481efc573b901fb7cb605e7505a34568d0e99cb7265bd1645ee2dac314e5d7fecd9375d46498325f03b6aefe444253f6a08431a3ffd9842838ba99"},
	SumTestCase{SKEIN1024_MAX - 7, -1, "beep boop", "d99c4d68d72f3800f3db6ccc89d0c47f2dad229e5dc64dd3cc8ab96105682a6afed0297bf778601edbf7a2e46a3a34cff773213da9de17357d81b1c6d2ba302174f2f0effba12f9a0824c792e327c092b6c975b458c8f491380a75071b6747ed6f9662723ecd09c671172febf56db816d534f2ae38bb763615"},
	SumTestCase{SKEIN1024_MAX - 6, -1, "beep boop", "2f2831ff120ef567f87e769729d1145d0059c0300c35ac77a6c40154c9e6e1c6dadd16797b459569bbd82f7707444e1c3383ed005941300a9b996e3ea21beac764edc975f440bff7234e9693ab1fb1c62620cb6868c8ca64bd5a491e6ab6dcc4e8a860a266c1eaf9c0082a0ea3c819489af299fe66f5982d774a"},
	SumTestCase{SKEIN1024_MAX - 5, -1, "beep boop", "c7a14a557824c189728db790b2a9f4b04c95f9cec37c7412946a8b7d578e2e80f307706617854910bf0c76d8908991478c43bf270a1f149fa59f62e8b629bc742d76f09456e4b89601f2f9ba8e5201c6badceb3a240ab9977a9f93a9011700db0564bf3cd437bb5baa5b670e50b16cc2f31a8de8e683ff9b9c70c3"},
	SumTestCase{SKEIN1024_MAX - 4, -1, "beep boop", "3e6d5a295fac29ab92ebb34ff6d61bb0bc2bd900724a8f4ccff981723f4f4cd6e3dd9469f9ec96c0aca1133541221eacb61490b3eaec3de985f2b773cab5cd98fd0a4cd36ef6ccf3a0e395469f51339770cb6d9edc69abe1e59e91dcf75422d2537aa7d0ee83794f56fab1dd45d26baade97005a103109f81e6db67d"},
	SumTestCase{SKEIN1024_MAX - 3, -1, "beep boop", "bcb96731d0a1fb40122f8872ae16f74300ec54bff65d391ae10f2de57f50b5fa38f7998058e51dc246f8da7e904bc74686bee3eeb6d48eb2245c21294af3fea5e36bae2fee5ea4fd785681c5a5b224b0b85774403ca010935f841ce8cb3e8a9dd78b2366fb71ee66c63e9a3ce11bc427c3d7642b3ba8ffad050a88553a"},
	SumTestCase{SKEIN1024_MAX - 2, -1, "beep boop", "e015846b4d493ddc813aeb6ed95a010952a95e553ce43182fc323c8c3de6c1526a985edb33e19ee38c94baf695ccfc90ca91b05b0975676380df0f0913e42e3fcff1be4b6e3a077657c565c0e27d50e49b4a713614456502ccb8c17073280f9d12e1f071e0d1dadb863941d7ab1fad29602dadeaa168ffb0e2abc41bcf83"},
	SumTestCase{SKEIN1024_MAX - 1, -1, "beep boop", "c0e5a0ab262ba3d58d6084d3ee7d1c911d7f5d653c421f4d186b1a0b3098c5409463f55c909270ca4c0acb89102ffc2642066bfbfe344fb49645d0d65f6307bcf2742c44347b01da3446175a795f8677e9463287307839a658de95fcdbd56cae1f7043b23e1bb64f4227acc9428df745c6e573468a087128afc6592c0979b2"},
	SumTestCase{SKEIN1024_MAX, -1, "beep boop", "789673e78a9719ce4f1555e83c6e8f060e1b0991e4e05eaa60c67bf6ff4af2b06a6a8ff640a1f1c39ecde50c77f26abef689707cd97eee6b0f5313f8524085e231d90735b3338fdae938c12e79586e70db03797fdf7ebb1dd7a1d8b4db793f615436cbf7913518f92d0ba69de873a4d544078e40099c1bcd6e97d80479598493"}
}

func TestSum(t *testing.T) {

	for _, tc := range sumTestCases {

		m1, err := FromHexString(tc.hex)
		if err != nil {
			t.Error(err)
			continue
		}

		m2, err := Sum([]byte(tc.input), tc.code, tc.length)
		if err != nil {
			t.Error(tc.code, "sum failed.", err)
			continue
		}

		if !bytes.Equal(m1, m2) {
			t.Error(tc.code, Codes[tc.code], "sum failed.", m1, m2)
			t.Error(hex.EncodeToString(m2))
		}

		s1 := m1.HexString()
		if s1 != tc.hex {
			t.Error("hex strings not the same")
		}

		s2 := m1.B58String()
		m3, err := FromB58String(s2)
		if err != nil {
			t.Error("failed to decode b58")
		} else if !bytes.Equal(m3, m1) {
			t.Error("b58 failing bytes")
		} else if s2 != m3.B58String() {
			t.Error("b58 failing string")
		}
	}
}

func TestBlakeMissing(t *testing.T) {
	data := []byte("abc")

	_, err := Sum(data, BLAKE2B_MAX-2, -1)
	if err == nil {
		t.Error("blake2b-496 shouldn't be supported")
	}

	_, err = Sum(data, BLAKE2S_MAX-2, -1)
	if err == nil {
		t.Error("blake2s-240 shouldn't be supported")
	}
}

func BenchmarkSum(b *testing.B) {
	tc := sumTestCases[0]
	for i := 0; i < b.N; i++ {
		Sum([]byte(tc.input), tc.code, tc.length)
	}
}
