package collect

import (
	"net/http"
	"github.com/PuerkitoBio/goquery"
	pb "crawler-apkmirror/protos"
	"crawler-apkmirror/config"
	"github.com/golang/glog"
	"strings"
	"fmt"
)

/*
    解析HTML页面
 */
func ParseAllVerHtml(resp *http.Response) (reply *pb.ApkDetailReply, err error) {

	var apkName, apkIcon, version, apkUpload, apksize, downloads, variantLink string
	doc, err := goquery.NewDocumentFromResponse(resp)

	if err != nil {
		glog.V(0).Infoln("NewDocumentFromResponse", err)
		return reply, err
	}
	//fmt.Println("ParseApkDetail doc",doc.Text()) widgetHeader
	doc.Find("#primary .listWidget ").Each(func(i int, selection *goquery.Selection) {
		text := selection.Find(".widgetHeader").Text()

		if text == config.ALLVERSIONS {

			apkName = selection.Find(".fontBlack").Eq(0).Text()
			apkIcon, _ = selection.Find(".table-cell img").Eq(0).Attr("src")
			version = selection.Find(".infoslide-value").Eq(0).Text()
			apkUpload = selection.Find(".infoslide-value").Eq(1).Text()
			apksize = selection.Find(".infoSlide .infoslide-value").Eq(2).Text()
			downloads = selection.Find(".infoSlide .infoslide-value").Eq(3).Text()
			//拿到apk不同版本链接
			variantLink, _ = selection.Find(".iconsBox .downloadLink ").Eq(0).Attr("href")

		}

	})

	//glog.V(0).Infof("ParseAllVerHtml variantLink %v\n", variantLink)

	apkLink, err := GetAllVerLink(variantLink)
	//glog.V(0).Infof("GetAllVerLink apkLink %v\n", apkLink)
	if err != nil {
		glog.V(0).Infoln("GetAllVerLink", err)
		return reply, err
	}

	return &pb.ApkDetailReply{
		ApkName:    apkName,
		ApkIcon:    config.MIRRPRURL + apkIcon,
		ApkVersion: version,
		ApkUpload:  apkUpload,
		ApkSize:    apksize,
		Downloads:  downloads,
		ApkLink:    apkLink,
	}, nil
}

func ParseLatestVerHtml(resp *http.Response) (reply []*pb.DownloadLink, err error) {

	variantDoc, err := goquery.NewDocumentFromResponse(resp)
	//fmt.Println("resp.StatusCode===", resp.StatusCode)
	if err != nil {
		glog.V(0).Infoln("variantDoc", err)
		return reply, err
	}

	variantDoc.Find(".listWidget .topmargin .headerFont .rowheight").Each(func(i int, selection *goquery.Selection) {

		downlaodUrl, _ := selection.Find("a").Attr("href")
		if downlaodUrl == "" {
			glog.V(0).Infof("variantDoc downlaodUrl is null")
			return

		}
		fmt.Println("downlaodUrl", downlaodUrl)
		downLink, err := GetLatestVerLink(config.MIRRPRURL + downlaodUrl)
		if err != nil {
			glog.Infoln("GetLatestVerLink", err)
			return
		}

		reply = append(reply, downLink)

	})
	//fmt.Println("reply====", reply)

	if len(reply) == 0 {

		var versionCode string
		downLink := &pb.DownloadLink{}
		href, _ := variantDoc.Find(".downloadButton").Attr("href")
		variantName := variantDoc.Find(".siteTitleBar .app-title").Text()
		//fmt.Println("href===before===", href)
		if !strings.Contains(href, "download.php?id=") || strings.HasSuffix(href, "download/") {
			href, _ = GetDownloadLink(config.MIRRPRURL + href)

		}
		//fmt.Println("href===after===", href)

		text := variantDoc.Find(".appspec-value").Text()
		//fmt.Println("text==========", text)
		if text != "" {
			versionCode = text[strings.Index(text, "(")+1:strings.Index(text, ")")]

		}
		//fmt.Println("versionCode===",versionCode)

		downLink.DownloadUrl = config.MIRRPRURL + href
		downLink.VariantName = variantName
		downLink.VersionCode = versionCode
		reply = append(reply, downLink)
	}
	//fmt.Println("reply==", reply)
	return reply, err
}

func ParseDownLoadHtml(resp *http.Response) (reply *pb.DownloadLink, err error) {
	var versionCode string
	downloadDoc, err := goquery.NewDocumentFromResponse(resp)
	if err != nil {
		glog.V(0).Infoln("downloadDoc", err)
		return reply, err
	}

	text := downloadDoc.Find(".appspec-value").Text()
	//fmt.Println("text==========", text)
	if text != "" {
		versionCode = text[strings.Index(text, "(")+1:strings.Index(text, ")")]

	}
	//fmt.Println("versionCode===",versionCode)
	variantName, _ := downloadDoc.Find(".siteTitleBar .app-title").Attr("title")
	downlaodUrl, _ := downloadDoc.Find(".noPadding .downloadButton").Attr("href")
	//glog.V(0).Infof("VariantName %s, downlaodUrl %s\n", variantName, downlaodUrl)

	//特殊情况下拿到的url不是下载链接，判断并重新抓取
	if !strings.Contains(downlaodUrl, "download.php?id=") || strings.HasSuffix(downlaodUrl, "download/") {

		downlaodUrl, _ = GetDownloadLink(config.MIRRPRURL + downlaodUrl)

	}

	return &pb.DownloadLink{
		VariantName: variantName,
		DownloadUrl: config.MIRRPRURL + downlaodUrl,
		VersionCode: versionCode,
	}, nil
}
