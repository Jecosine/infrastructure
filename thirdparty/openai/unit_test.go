package openai

import (
	"bytes"
	"github.com/alioth-center/infrastructure/logger"
	"github.com/alioth-center/infrastructure/network/http"
	"io"
	h "net/http"
	"os"
	"strings"
	"testing"
)

func TestOpenAiClient(t *testing.T) {
	// uses real openai endpoint to test, because mocking the endpoint cannot find issues in the client
	var client Client
	apiKey, baseUrl := os.Getenv("OPENAI_API_KEY"), os.Getenv("OPENAI_BASE_URL")
	if apiKey == "" || baseUrl == "" {
		t.Log("OPENAI_API_KEY or OPENAI_BASE_URL is not set, mock it, but it is not recommended")
		client = initMockingClient(t)
	} else {
		client = NewClient(Config{ApiKey: apiKey, BaseUrl: baseUrl}, logger.Default())
	}

	t.Run("CompleteChat", func(t *testing.T) {
		response, err := client.CompleteChat(CompleteChatRequest{
			Body: CompleteChatRequestBody{
				Model: "gpt-4o",
				Messages: []ChatMessageObject{
					{
						Role:    ChatRoleEnumSystem,
						Content: "now testing api is working, please echo any input",
					},
					{
						Role:    ChatRoleEnumUser,
						Content: "testing",
					},
				},
				N: 1,
			},
		})

		if err != nil {
			t.Error(err)
		}

		if len(response.Choices) == 0 || !strings.Contains(response.Choices[0].Message.Content, "testing") {
			t.Error("response is not as expected")
		}
	})

	t.Run("Embedding", func(t *testing.T) {
		response, err := client.Embedding(EmbeddingRequest{
			Body: EmbeddingRequestBody{
				Input: "Hello, world!",
			},
		})

		if err != nil {
			t.Error(err)
		}

		if len(response.Data) == 0 {
			t.Error("response is not as expected")
		}
	})
}

func initMockingClient(t *testing.T) Client {
	t.Helper()

	var (
		chatCompletionsOpts = &http.MockOptions{
			Trigger: func(req *h.Request) bool {
				if strings.HasSuffix(req.URL.String(), "/v1/chat/completions") {
					return true
				}

				return false
			},
			Handler: func(req *h.Request) *h.Response {
				return &h.Response{
					Status:     "ok",
					StatusCode: 200,
					Body:       io.NopCloser(bytes.NewBufferString(`{"choices":[{"finish_reason":"stop","index":0,"logprobs":null,"message":{"content":"testing","role":"assistant"}}],"created":1719546132,"id":"chatcmpl-9ewmyvT7IZL4J3uMRV1d1TIiYyVIQ","model":"gpt-4o-2024-05-13","object":"chat.completion","system_fingerprint":"fp_ce0793330f","usage":{"completion_tokens":1,"prompt_tokens":22,"total_tokens":23}}}`)),
				}
			},
		}
		embeddingOpts = &http.MockOptions{
			Trigger: func(req *h.Request) bool {
				if strings.HasSuffix(req.URL.String(), "/v1/embeddings") {
					return true
				}

				return false
			},
			Handler: func(req *h.Request) *h.Response {
				return &h.Response{
					Status:     "ok",
					StatusCode: 200,
					Body:       io.NopCloser(bytes.NewBufferString(`{"data":[{"embedding":[0.0014316267,0.00338094,-0.013127507,-0.033413146,-0.009503701,0.0049116304,-0.015389988,0.0018742167,-0.0029687083,-0.024963992,0.029910775,0.007183699,-0.016872745,-0.01797203,0.010392076,-0.002946339,0.025015121,-0.015236599,0.011408276,0.011012022,-0.008289374,-0.0018550432,0.017051697,0.005742485,-0.014367398,-0.0075224317,0.0035854583,-0.015799025,0.03750351,-0.026088841,0.010027778,-0.0066021,-0.004834936,-0.013868884,0.011849268,-0.019173574,0.005080997,-0.0115361,0.019467568,-0.011855659,0.004755046,0.005707334,0.003166835,-0.006157913,-0.026715178,0.008340504,0.0054772506,-0.009164968,-0.0062378026,0.020439029,0.020579634,-0.012053786,-0.013510978,0.007861165,-0.018790102,-0.00042141916,-0.03837271,0.024120355,0.03101006,-0.022228563,0.019940516,0.008091248,-0.02312333,0.009951085,-0.010072517,0.0049723466,0.014533568,0.013945579,-0.02221578,0.024171485,0.020375118,0.0018310762,-0.0055731186,-0.004700721,0.02018338,0.00029020003,-0.020605199,0.005892678,0.00039126072,-0.008941276,0.028811488,-0.036148578,-0.013344807,0.014546351,0.022573687,0.0036557612,-0.005649813,0.022918811,-0.016821615,-0.015594506,-0.010040561,0.006004524,0.015479465,0.014047838,-0.024056442,0.012910206,0.0028472757,0.031240141,-0.00026703195,-0.035816234,-0.013421501,-0.0036110228,-0.010449598,-0.0128526855,-0.019454785,0.0071133957,-0.011274061,-0.0021042996,0.027686639,0.001525098,-0.02323837,0.008922103,-0.0005728105,-0.033489842,-0.0023759252,-0.00864089,-0.00093870616,0.010660507,0.0013149875,-0.020375118,0.018802883,0.02135936,0.00633367,-0.023672972,0.0020419855,0.0038283234,-0.034998164,-0.021768397,-0.0076055173,-0.01271847,0.025296334,-0.0019764758,-0.0042245775,0.000118336895,-0.033464275,0.027482122,-0.011325191,0.019608174,-0.021934567,-0.009574004,0.00086680526,0.028018981,0.011951528,-0.0020212142,-0.0071773077,0.020899195,0.0038922352,-0.013715496,0.013280895,-0.0067427065,0.0058703087,0.018802883,0.011127064,0.0039721252,0.012545908,0.035074856,-0.0018710211,0.020681893,0.004416313,-0.017511863,0.003815541,-0.008755932,0.006391191,-0.025628677,0.017345693,0.019173574,0.019697651,0.011274061,-0.009963867,-0.020400682,-0.012066568,0.002866449,-0.043127757,0.022573687,-0.014776434,0.015057646,0.0038666707,0.01625919,-0.014290703,-0.02934835,-0.027993416,0.008883756,0.0037324557,0.029757386,0.009005188,-0.00013880868,0.012775991,-0.011663924,0.0048764786,-0.013958361,0.008685629,0.037708025,0.024772257,0.017524647,-0.6839086,-0.012028221,0.01942922,0.01916079,0.0061163697,0.009867999,0.0076438645,0.037759155,0.00028380883,0.030243115,-0.021001454,0.02296994,-0.014239574,-0.0066851857,-0.0075735613,-0.014993734,0.004020059,-0.015581723,-0.004860501,0.006998354,-0.0127248615,0.02914383,-0.023788013,-0.016757702,0.011785356,0.008557805,0.0065765358,-0.015338859,-0.014431309,0.03267177,-0.0040424285,0.013191419,0.005624248,0.0059981328,0.049723465,0.016080236,0.0023423715,0.0012958139,0.007963424,0.029859645,-0.019978862,-0.01985104,0.010890589,-0.016463708,0.00890932,0.0048860656,0.0030709673,-0.0056338347,0.0038762575,0.0072156545,0.016118584,-0.0005923835,0.015044863,0.019084096,-0.0008340504,0.007758906,0.015914066,0.014533568,0.00008313541,0.005643422,0.0014076598,0.02130823,-0.0016760898,0.006138739,-0.019505914,0.018738972,-0.018393848,0.015249382,-0.003815541,-0.012200784,-0.012699297,-0.006928051,-0.0150192985,-0.018342718,0.006154717,0.025922671,0.02580763,-0.003428874,-0.0053686006,0.017562993,0.017371258,0.017192304,-0.0064263428,-0.023647407,0.016463708,-0.023327848,-0.030779976,-0.017844206,0.0019652913,0.0043492056,0.03745238,0.0070111365,-0.010315382,-0.030115291,0.024286525,0.0061738905,-0.012603429,-0.0021841896,0.017793076,-0.008986015,0.01041125,0.00086440856,0.0076822117,-0.0012223152,0.016233625,0.012513952,-0.001074519,0.041926213,0.021282665,-0.025654241,-0.005279124,0.005710529,-0.037708025,0.015159905,0.011708662,-0.024554957,0.001374905,0.022036826,0.029731821,-0.014750869,0.02296994,-0.008449155,0.009855216,0.007854774,0.002503749,0.015159905,-0.013421501,0.006432734,-0.016463708,0.005576314,0.018278806,0.007560779,0.0009866401,-0.0034895902,0.004889261,0.006883313,0.015249382,-0.011357146,0.0033553753,0.02102702,-0.007841991,0.0035407199,0.00032435294,-0.001714437,0.009887173,-0.018969055,-0.0017160348,0.017959246,0.005838353,-0.0060844137,0.014623045,-0.008839017,-0.028428018,0.028658101,-0.011510535,-0.00214744,-0.007720559,-0.025833193,-0.026331708,-0.031035624,0.0100853,0.01738404,-0.011101499,0.003908213,-0.011740617,-0.02720091,-0.020669112,0.023059417,0.004659178,-0.022867681,0.005371796,-0.011414668,-0.024439914,-0.0038666707,-0.0026747135,0.010845851,-0.012481996,-0.022842117,-0.0013077975,-0.020247294,0.0045856796,-0.007969815,-0.0022465037,0.018202111,0.027073085,-0.002607606,0.0021746028,0.022931593,-0.024618868,-0.010769157,0.01706448,0.011977092,-0.019352527,-0.004438682,0.0019732802,-0.01803594,-0.0027194517,-0.0017495885,0.008877364,0.019672086,0.020681893,0.007861165,0.020311205,-0.026152754,-0.020042775,-0.026919696,-0.0007585545,-0.031086754,0.0051385174,0.016540403,0.01142745,0.0012846293,-0.0003353378,0.016629878,0.017997595,0.024273744,-0.0012966129,-0.022599252,-0.022675946,0.0059054606,0.0038858443,0.014303486,-0.016169714,-0.006589318,0.008628108,0.03149579,0.007976206,0.04046902,0.010743592,-0.02704752,-0.019224703,0.017754728,0.011612794,0.014584698,-0.008033727,-0.0023407738,0.0019589001,-0.0130252475,0.035356067,0.012846294,-0.0075735613,0.013600455,0.026945261,-0.017767511,-0.0031652374,0.030243115,0.020630764,-0.0066468385,0.0034544386,0.028504713,-0.010909763,-0.0044258996,-0.001351737,-0.014546351,0.014750869,-0.030831106,0.008704802,-0.01518547,0.021333795,0.032595076,0.020988671,0.008487502,0.0033713533,-0.006541384,0.004173448,-0.0130891595,-0.0054708594,-0.014073403,0.0053014928,-0.0007781275,-0.022509774,-0.011446623,0.027865592,-0.03783585,0.005950199,0.0032850723,-0.0047933934,-0.010289817,-0.0073051313,0.0007437749,-0.008366069,-0.026868567,0.0014460069,0.006525406,0.004834936,-0.012993291,-0.029987467,-0.009695437,-0.0054389033,-0.0035215463,-0.000058619204,0.009049927,0.0074009993,-0.026433965,0.0118684415,-0.0038027586,0.02548807,0.0031061189,-0.0012614613,-0.015977977,0.014904258,0.022739857,-0.0066021,-0.018176548,0.0053526224,0.00034891907,-0.006755489,-0.022535339,0.006960007,-0.015402771,0.006193064,0.0029271655,-0.00073738367,0.00402645,0.010877807,0.0011400287,-0.00089556567,0.023890272,0.029757386,0.0064678853,0.0068130097,-0.001390084,-0.021921786,-0.008545022,0.076592036,0.0011392297,0.0025964214,0.015926847,-0.015914066,-0.008768714,-0.02607606,-0.021781178,-0.02119319,-0.0113315815,-0.0012854283,0.0023200023,-0.017754728,0.008187116,0.008116812,0.008225462,0.004406726,-0.015632853,-0.0076758205,0.010021388,-0.008832626,0.009574004,-0.002703474,0.010954501,0.024938427,0.0008284581,0.006582927,0.012616211,-0.008129595,-0.009759349,0.0044993986,-0.010053343,0.018623931,0.0017799467,0.013587672,0.014277921,-0.00922888,0.029220525,0.010647724,-0.010845851,0.0067171417,0.0032531163,-0.009088274,-0.0026747135,0.010008605,0.002361545,-0.009113838,0.028146805,-0.0025500853,-0.027252039,0.021972915,0.0015035276,-0.008097639,-0.014098967,-0.014035055,0.015581723,0.004055211,-0.016016325,-0.014980951,-0.018202111,-0.006966398,-0.018419413,-0.016744921,-0.009637916,-0.02199848,-0.026919696,-0.022893246,-0.0046751564,-0.011312408,-0.0059310254,-0.010820286,-0.045479715,-0.016617097,-0.0038602795,0.01926305,0.0018342718,0.005742485,-0.006033284,0.0050682146,0.0020947128,-0.014009491,-0.027252039,0.012974118,-0.0106668975,0.014303486,0.0039465604,-0.0042309687,0.009286401,-0.022177434,0.031521354,0.0021777984,0.024299309,0.022190215,0.0017799467,0.007816426,0.007631082,-0.0006866536,-0.0012303042,-0.01824046,-0.0005188848,-0.005592292,-0.0086089345,-0.0011815714,-0.010149212,-0.00042301696,0.008698411,0.010992848,0.024184266,-0.01599076,-0.00006795634,0.034819208,-0.022906028,0.018048724,0.002433446,0.0026044103,0.0075096493,0.027942287,0.031674743,0.00046855418,-0.03218604,-0.0021122887,-0.022548122,0.0064359293,0.015594506,0.006787445,-0.00028820278,-0.0022768618,-0.015786242,-0.010762766,-0.0007697391,-0.006550971,0.016463708,-0.0025229226,-0.0069152685,-0.005617857,-0.0083085485,-0.020630764,0.010954501,-0.0043364232,-0.015568941,-0.017473517,-0.006525406,-0.002281655,-0.014405744,-0.04448269,-0.04266759,-0.0073115225,-0.0022848507,-0.007656647,0.018560018,-0.0027673857,0.00070023484,-0.018956272,-0.014520786,-0.0015490649,-0.020669112,-0.003364962,-0.0038570839,0.02699639,0.025027905,0.023327848,0.015032081,0.016438143,0.0002868047,-0.0007725352,-0.00987439,0.009248054,-0.005080997,-0.011957918,0.0076694293,0.017729163,0.018150983,-0.013677149,0.017320128,0.0013046019,-0.008334113,-0.015645636,0.011069543,-0.01760134,-0.017703598,-0.028402453,-0.0018119026,0.00410634,0.004141492,0.0038443014,-0.01159362,0.016438143,0.012366954,0.019672086,0.0037963674,0.033771053,-0.024772257,0.024299309,0.009113838,0.013242547,-0.016987786,-0.0027833637,-0.012143263,0.022918811,0.004256533,0.017767511,0.023315065,-0.0068385745,-0.009331139,0.0012990095,0.027661074,-0.005675378,0.0028328954,0.0060492624,-0.0060556536,0.002195374,-0.021448838,-0.00024865728,-0.013293677,0.0043044672,-0.015862936,-0.019288614,0.01207296,-0.017575776,-0.03520268,-0.019774346,0.0006618877,0.04599101,-0.0040711886,0.009299183,0.022701511,0.0026171927,0.001209533,0.011555273,-0.016859962,-0.004445073,0.03827045,0.007822818,-0.010833069,-0.013069985,-0.0005228793,0.0055507496,-0.013421501,-0.01824046,-0.0009283205,0.023519583,0.0024845756,-0.016502054,0.006026893,-0.028402453,0.021806743,-0.0023216002,-0.021372143,-0.020272858,0.0071006133,-0.023404542,0.016361449,-0.02033677,0.033413146,0.023200024,-0.014725304,-0.006515819,0.0043172496,0.0008316537,-0.009772131,-0.0007669429,0.025424158,-0.0035183507,0.00043220428,0.013050812,-0.0016457317,-0.016834397,-0.008251027,-0.007055875,0.023800796,-0.023903055,0.0074137817,-0.00557951,0.009318356,-0.000014417628,0.014520786,-0.014827563,-0.015837371,0.003396918,-0.0005416534,0.033106368,0.0027194517,0.00089316897,-0.018393848,-0.013268112,-0.027252039,-0.051717516,-0.0006359235,-0.005704138,-0.03407783,-0.013587672,-0.0038251278,0.009599569,0.018675061,0.0179209,-0.009995823,0.0136260195,0.023864707,-0.015926847,-0.005656204,0.02720091,0.014827563,-0.028530277,0.018560018,-0.003515155,-0.020630764,0.0012958139,-0.007912295,-0.010762766,0.0071709165,0.011887616,0.008858191,-0.009976649,0.008493893,0.029169396,0.028811488,0.018623931,-0.0074329553,-0.009887173,0.047831673,-0.009522875,-0.000105554514,-0.0049755424,-0.012462823,0.02431209,-0.007420173,-0.016681008,0.008027336,-0.03740125,-0.004953173,-0.005624248,-0.0062218243,0.009612352,-0.025296334,-0.02233082,-0.007094222,-0.017895335,0.0024765865,0.008660064,-0.016617097,-0.0030214356,0.017102826,0.019084096,-0.0042693154,-0.0005512402,-0.011274061,-0.017665252,-0.014584698,-0.022087956,-0.01513434,-0.0038506927,0.041414917,0.012475605,-0.011229322,-0.009746566,0.002275264,-0.026715178,-0.033464275,-0.0015786241,-0.0060556536,0.0014715717,0.013012465,-0.0030549893,0.012098525,0.004684743,0.00087719096,-0.013306459,-0.016629878,0.021372143,-0.003120499,0.018010376,0.005224799,-0.03037094,-0.0017575775,0.001990856,0.014789216,-0.011299626,0.012577864,-0.00450579,-0.012475605,-0.01110789,-0.0107308095,0.0003668943,0.008934885,0.021985697,-0.012526735,-0.0051864516,0.00036789294,-0.007631082,0.005854331,-0.009484528,0.014533568,-0.018393848,0.01325533,-0.021640573,-0.005256755,-0.0072156545,-0.017128391,0.0069152685,-0.017959246,0.00788673,0.024567738,0.0009099458,0.009209706,-0.00450579,0.00011933552,-0.006397582,0.009318356,0.00021470408,-0.013012465,0.0026587355,0.0034863946,-0.018790102,-0.025590328,0.005189647,-0.036020752,-0.008883756,-0.0064167557,0.034052264,0.015952412,-0.020515723,-0.008717584,-0.01990217,0.00014330249,0.01062216,-0.017984811,0.0078100353,-0.0075735613,0.0063879956,0.025321899,-0.016361449,0.0024622064,0.01760134,0.013510978,-0.0065126237,0.026945261,0.23294613,-0.0118045295,0.010487945,0.024350438,0.013894449,0.017447952,0.010424033,-0.00096187426,0.0015618473,0.014418527,-0.009721002,-0.004818958,-0.0009155381,0.0064263428,0.0022001674,0.001619368,-0.018227676,-0.015581723,-0.029322784,-0.028428018,-0.0052184076,-0.030447634,-0.019672086,-0.017933682,0.017409604,-0.009375878,-0.006550971,0.0027418209,0.033438712,0.0062122378,-0.006889704,-0.014814781,0.010954501,0.0065957094,-0.006691577,-0.02403088,0.03246725,0.00560827,0.03788698,0.0074904757,-0.0103025995,-0.01706448,-0.005049041,-0.014763651,0.011184584,0.009663481,0.0022704706,-0.026229449,0.008519458,0.03456356,0.00073738367,0.009695437,0.0125331255,0.03364323,0.00024106774,-0.013894449,0.011747009,0.013702714,0.0021058975,-0.0039465604,-0.00008877764,0.027916722,-0.0043044672,0.02812124,-0.022151869,0.0073690433,-0.031035624,-0.0010601388,0.008449155,-0.009976649,-0.015479465,-0.00480298,0.01325533,-0.0058671134,-0.029092701,-0.0039178003,0.014137315,-0.004400335,0.025475288,0.017793076,-0.0030869453,0.0055443584,-0.0056593996,-0.018751755,-0.03236499,-0.028862618,0.0048796744,-0.0027194517,-0.015517812,-0.023915837,-0.009784914,-0.01148497,-0.000045537236,-0.0017767511,0.02699639,0.024158701,-0.0021138864,0.014559133,-0.013894449,0.0082446365,-0.007496867,0.003002262,-0.00552838,0.012744035,-0.004556919,0.0029367523,-0.009676264,-0.00009082482,-0.001050552,-0.023864707,-0.029655127,-0.0066724033,0.004384357,-0.00203879,-0.008864582,0.0018294784,-0.00043899493,0.0057872236,-0.005291906,-0.011881224,-0.0022624817,-0.012156045,-0.013677149,-0.0004589674,-0.005039454,-0.021870656,-0.0071389605,-0.0040168636,0.013408719,-0.029450608,0.015914066,-0.0009195326,0.027763333,0.014316268,-0.006896095,0.010973675,0.014022273,-0.005908656,-0.016182495,-0.0009738577,0.002465402,-0.02248421,0.016271973,-0.0066787945,-0.0143546155,-0.010769157,0.035074856,-0.026357273,-0.00014040648,0.0067363153,-0.035330504,-0.012079351,-0.008564196,-0.011478579,0.026894132,-0.01689831,-0.022637598,-0.026638484,-0.0023439694,0.002946339,-0.026178319,0.012271087,0.0362764,-0.0015946021,-0.026178319,-0.013140289,-0.16392127,0.015709547,0.016284754,-0.020720242,0.027226474,0.010686072,0.030728847,0.0026187906,-0.0057520717,0.021282665,0.008187116,-0.018764537,-0.014073403,-0.008001771,0.0042852936,-0.009273618,-0.0065349927,0.017524647,0.03223717,0.027788898,0.026970826,-0.019339744,0.010430424,-0.0032323448,0.01164475,-0.009599569,-0.003815541,0.037273426,-0.016859962,-0.008442763,-0.016617097,-0.022931593,0.022867681,0.011747009,0.01760134,0.0003069769,0.012846294,-0.044303738,-0.015977977,0.02011947,0.025398593,0.011580838,0.02382636,-0.016719356,-0.015364423,-0.008564196,0.020835282,0.01921192,0.024836168,-0.009427006,0.010155602,-0.004969151,0.0037931718,-0.01738404,0.02763551,0.028683666,-0.00097146106,0.01786977,0.00028281022,-0.005055432,-0.00023227985,-0.016808832,-0.0051800604,0.0025564765,0.011088717,-0.007931468,0.004841327,0.013971143,-0.0005049041,0.012865468,-0.00858976,-0.024580522,-0.010909763,-0.031240141,0.0045856796,0.0060492624,-0.03374549,0.018713407,-0.014201227,-0.026357273,-0.010059735,0.027865592,0.00007319911,0.0035854583,-0.018317154,0.00003193099,0.0051800604,0.018470543,-0.031291272,-0.04067354,0.0143546155,-0.0125970375,-0.007963424,-0.02135936,0.00031057195,0.017626904,0.018483324,0.010826678,-0.017205086,-0.013216983,0.002013225,0.00847472,-0.0183555,0.01379219,0.03970208,-0.005237581,0.008033727,0.004732677,0.024746692,-0.026165536,-0.032748464,0.007023919,0.015428335,0.027226474,0.013715496,0.023813577,-0.0030853474,-0.017780293,0.003451243,-0.015032081,0.053225838,-0.0049403906,-0.043255582,0.0005392567,0.0030454025,-0.01051351,-0.08732923,-0.026306143,-0.0004873283,0.009567613,0.0066532297,0.04818958,0.0062921275,0.0055475538,-0.00016966615,0.009880781,0.004384357,-0.03021755,-0.025155729,-0.011970701,0.042642027,-0.032390557,0.01303803,0.0012766405,-0.02812124,0.013843319,-0.0056306394,-0.00044738338,0.00088278326,-0.008295766,-0.017486298,0.012616211,-0.020758588,0.03095893,0.0032371383,-0.012341389,-0.021474402,-0.027456557,-0.012188002,-0.008084857,0.0037324557,0.011772574,-0.0407758,0.031904824,0.012686514,-0.04013668,0.02607606,0.011663924,-0.006365626,-0.0644232,-0.003959343,0.018521672,-0.001722426,0.035713974,0.011657532,-0.00998304,-0.02651066,-0.025424158,-0.026229449,0.0040008854,0.01899462,-0.006282541,0.008826234,0.030805541,0.003908213,0.011005631,0.011075934,0.004502594,-0.020221729,0.015389988,-0.0074713025,0.018163765,-0.0154539,0.024848951,0.005838353,0.0013085963,-0.019991646,0.021857873,0.007420173,0.014456874,-0.033259757,-0.006359235,-0.019697651,-0.010347338,0.010865024,-0.00089636457,-0.034435738,-0.01694944,0.018010376,-0.00080009724,0.017473517,0.010257862,-0.013549325,0.002535705,-0.00821268,-0.03438461,-0.028709231,0.0022401125,0.011261279,-0.0066148825,0.010046952,0.0069536157,-0.006582927,-0.00014340234,0.018585583,0.012884641,-0.03246725,0.0029958708,-0.065445796,0.02978295,-0.014303486,-0.031777002,-0.0006107582,-0.014853128,0.008915711,-0.015095993,0.0006822596,-0.012871859,-0.014571915,-0.0064039733,0.009305574,0.014942605,-0.030984495,0.00606524,0.0078100353,0.012833511,0.009989431,0.006240998,0.010545465,-0.036199708,-0.008807061,-0.0049979114,-0.024248179,0.0078356005,-0.03369436,-0.00230722,-0.0035247419,-0.010455988,0.013894449,-0.041517176,0.009931911,0.021819526,0.010161994,-0.008711194,0.0046176356,0.022586469,-0.000845235,0.024529392,-0.028428018,-0.02763551,0.010219514,-0.016783267,-0.006902486,0.0043747704,-0.004112731,0.001358128,0.034103394,0.021448838,-0.009241662,0.03423122,-0.03218604,-0.015926847,-0.015083211,-0.014201227,-0.005493229,0.004314054,-0.0024765865,-0.038347147,0.029220525,0.006627665,0.027175345,-0.005831962,0.006697968,0.0018997815,-0.023532365,0.0012854283,-0.018061506,-0.022522558,-0.009721002,-0.012916597,0.0050362586,-0.0016377427,0.018496107,0.020720242,-0.0071325693,0.013242547,-0.043920264,0.03637866,0.014329051,0.016911091,-0.029731821,-0.0013765028,0.03970208,0.012705687,-0.008206289,-0.010845851,-0.0019349331,-0.0064998413,-0.028095676,-0.017102826,0.0074457377,0.027507687,0.003547111,0.010660507,0.021627791,0.01105037,0.02130823,0.02419705,-0.005413339,0.016859962,-0.020579634,-0.0067682713,-0.016617097,0.010245079,-0.0065349927,-0.04062241,0.0077908617,0.010449598,0.010673289,0.035458326,-0.0007433754,0.011823704,-0.026868567,-0.010264253,0.02419705,0.0049947156,-0.02607606,0.043562356,-0.013830537,0.0057840277,0.022867681,-0.021116495,0.004170252,-0.00046815473,-0.012334999,-0.011299626,0.012002657,-0.001192756,0.018150983,-0.004914826,-0.023327848,-0.004259729,-0.00007909099,0.020707458,-0.0071773077,0.0257565,-0.014162879,0.05215212,-0.007905903,-0.0019189551,0.022893246,0.022394734,0.014137315,0.018636713,0.01582459,-0.009995823,-0.015952412,-0.006634056,0.010270644,0.0066021,-0.027763333,-0.025743717,0.0022704706,-0.02189622,-0.0009794501,-0.015428335,0.018738972,0.018636713,0.013498195,0.028888183,0.0061994554,-0.024708344,-0.014201227,0.0033777445,0.0074904757,-0.007234828,-0.023903055,0.016182495,0.00922888,-0.010507118,-0.012104915,-0.000057770376,0.0025772478,-0.013562107,-0.0036749349,-0.00086760416,0.0059374166,0.024465479,0.004445073,-0.025539199,-0.026587354,0.004090362,0.005343036,-0.03095893,-0.0024606085,-0.019122444],"index":0,"object":"embedding"}],"model":"text-embedding-ada-002","object":"list","usage":{"completion_tokens":0,"prompt_tokens":4,"total_tokens":4}}}`)),
				}
			},
		}
	)

	client := http.NewMockClientWithLogger(logger.Default(), chatCompletionsOpts, embeddingOpts)
	return NewCustomClient(Config{}, client)
}