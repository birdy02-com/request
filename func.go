package request

import (
	"bytes"
	"fmt"
	"github.com/jpillora/go-tld"
	"github.com/yinheli/mahonia"
	"golang.org/x/net/html/charset"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// About IPv4

// GetRandomIP 生成一个随机的IP地址
func GetRandomIP() string {
	return fmt.Sprintf("101.%d.%d.%d", rand.Intn(255-1)+1, rand.Intn(255-1)+1, rand.Intn(255-1)+1)
}

// IsIpv4 检查一个IP字符串是否为IPv4地址
func IsIpv4(ip string) bool {
	if ip == "" {
		return true
	}
	parsedIP := net.ParseIP(ip)
	return parsedIP != nil && parsedIP.To4() != nil
}

// IsPrivateIP 判断 IP 地址是否为私有地址
func IsPrivateIP(ipStr string) bool {
	if IsIpv4(ipStr) {
		ip := net.ParseIP(ipStr)
		if ip.IsLoopback() {
			return true
		}
		return ip.IsPrivate()
	}
	return false
}

// About Domain

// RootDomain 检查一个字符串是否为域名并获取根域名
func RootDomain(domain string) string {
	if strings.Contains(domain, "@") {
		return ""
	}
	domSuffix := []string{".ac.ae", ".co.ae", ".net.ae", ".org.ae", ".sch.ae", ".cargo.aero", ".charter.aero", ".com.af", ".edu.af", ".gov.af", ".net.af", ".org.af", ".co.ag", ".com.ag", ".net.ag", ".nom.ag", ".org.ag", ".com.ai", ".net.ai", ".off.ai", ".org.ai", ".com.al", ".edu.al", ".net.al", ".org.al", ".co.am", ".com.am", ".net.am", ".north.am", ".org.am", ".radio.am", ".south.am", ".co.ao", ".it.ao", ".og.ao", ".pb.ao", ".com.ar", ".int.ar", ".net.ar", ".org.ar", ".co.at", ".or.at", ".asn.au", ".com.au", ".id.au", ".info.au", ".net.au", ".org.au", ".com.aw", ".biz.az", ".co.az", ".com.az", ".edu.az", ".gov.az", ".info.az", ".int.az", ".mil.az", ".name.az", ".net.az", ".org.az", ".pp.az", ".pro.az", ".co.ba", ".com.ba", ".co.bb", ".com.bb", ".net.bb", ".org.bb", ".ac.bd", ".com.bd", ".net.bd", ".org.bd", ".biz.bh", ".cc.bh", ".com.bh", ".edu.bh", ".me.bh", ".name.bh", ".net.bh", ".org.bh", ".co.bi", ".com.bi", ".edu.bi", ".info.bi", ".mo.bi", ".net.bi", ".or.bi", ".org.bi", ".auz.biz", ".com.bj", ".edu.bj", ".com.bm", ".net.bm", ".org.bm", ".com.bn", ".net.bn", ".org.bn", ".com.bo", ".net.bo", ".org.bo", ".tv.bo", ".abc.br", ".adm.br", ".adv.br", ".agr.br", ".am.br", ".aparecida.br", ".app.br", ".arq.br", ".art.br", ".ato.br", ".belem.br", ".bhz.br", ".bib.br", ".bio.br", ".blog.br", ".bmd.br", ".boavista.br", ".bsb.br", ".campinas.br", ".caxias.br", ".cim.br", ".cng.br", ".cnt.br", ".com.br", ".coop.br", ".curitiba.br", ".des.br", ".det.br", ".dev.br", ".ecn.br", ".eco.br", ".emp.br", ".enf.br", ".eng.br", ".esp.br", ".etc.br", ".eti.br", ".far.br", ".flog.br", ".floripa.br", ".fm.br", ".fnd.br", ".fortal.br", ".fot.br", ".foz.br", ".fst.br", ".g12.br", ".geo.br", ".ggf.br", ".gru.br", ".imb.br", ".ind.br", ".inf.br", ".jampa.br", ".jor.br", ".lel.br", ".logger.br", ".macapa.br", ".maceio.br", ".manaus.br", ".mat.br", ".med.br", ".mil.br", ".mus.br", ".natal.br", ".net.br", ".nom.br", ".not.br", ".ntr.br", ".odo.br", ".org.br", ".palmas.br", ".poa.br", ".ppg.br", ".pro.br", ".psc.br", ".psi.br", ".qsl.br", ".radio.br", ".rec.br", ".recife.br", ".rep.br", ".rio.br", ".salvador.br", ".seg.br", ".sjc.br", ".slg.br", ".srv.br", ".taxi.br", ".tec.br", ".teo.br", ".tmp.br", ".trd.br", ".tur.br", ".tv.br", ".vet.br", ".vix.br", ".vlog.br", ".wiki.br", ".zlg.br", ".com.bs", ".net.bs", ".org.bs", ".com.bt", ".org.bt", ".ac.bw", ".co.bw", ".net.bw", ".org.bw", ".com.by", ".minsk.by", ".net.by", ".co.bz", ".com.bz", ".net.bz", ".org.bz", ".za.bz", ".com.cd", ".net.cd", ".org.cd", ".ac.ci", ".co.ci", ".com.ci", ".ed.ci", ".edu.ci", ".go.ci", ".in.ci", ".int.ci", ".net.ci", ".nom.ci", ".or.ci", ".org.ci", ".biz.ck", ".co.ck", ".edu.ck", ".gen.ck", ".gov.ck", ".info.ck", ".net.ck", ".org.ck", ".co.cm", ".com.cm", ".net.cm", ".ac.cn", ".mil.cn", ".ah.cn", ".bj.cn", ".com.cn", ".edu.cn", ".gov.cn", ".cq.cn", ".fj.cn", ".gd.cn", ".gs.cn", ".gx.cn", ".gz.cn", ".ha.cn", ".hb.cn", ".he.cn", ".hi.cn", ".hk.cn", ".hl.cn", ".hn.cn", ".jl.cn", ".js.cn", ".jx.cn", ".ln.cn", ".mo.cn", ".net.cn", ".nm.cn", ".nx.cn", ".org.cn", ".qh.cn", ".sc.cn", ".sd.cn", ".sh.cn", ".sn.cn", ".sx.cn", ".tj.cn", ".tw.cn", ".xj.cn", ".xz.cn", ".yn.cn", ".zj.cn", ".com.co", ".net.co", ".nom.co", ".ae.com", ".africa.com", ".ar.com", ".br.com", ".cn.com", ".co.com", ".de.com", ".eu.com", ".gb.com", ".gr.com", ".hk.com", ".hu.com", ".it.com", ".jpn.com", ".kr.com", ".mex.com", ".no.com", ".nv.com", ".pty-ltd.com", ".qc.com", ".ru.com", ".sa.com", ".se.com", ".uk.com", ".us.com", ".uy.com", ".za.com", ".co.cr", ".ed.cr", ".fi.cr", ".go.cr", ".or.cr", ".sa.cr", ".com.cu", ".com.cv", ".edu.cv", ".int.cv", ".net.cv", ".nome.cv", ".org.cv", ".publ.cv", ".com.cw", ".net.cw", ".ac.cy", ".biz.cy", ".com.cy", ".ekloges.cy", ".ltd.cy", ".name.cy", ".net.cy", ".org.cy", ".parliament.cy", ".press.cy", ".pro.cy", ".tm.cy", ".co.cz", ".co.de", ".com.de", ".biz.dk", ".co.dk", ".co.dm", ".com.dm", ".net.dm", ".org.dm", ".art.do", ".com.do", ".net.do", ".org.do", ".sld.do", ".web.do", ".com.dz", ".com.ec", ".fin.ec", ".info.ec", ".med.ec", ".net.ec", ".org.ec", ".pro.ec", ".co.ee", ".com.ee", ".fie.ee", ".med.ee", ".pri.ee", ".com.eg", ".edu.eg", ".eun.eg", ".gov.eg", ".info.eg", ".name.eg", ".net.eg", ".org.eg", ".tv.eg", ".com.es", ".edu.es", ".gob.es", ".nom.es", ".org.es", ".biz.et", ".com.et", ".info.et", ".name.et", ".net.et", ".org.et", ".biz.fj", ".com.fj", ".info.fj", ".name.fj", ".net.fj", ".org.fj", ".pro.fj", ".co.fk", ".radio.fm", ".aeroport.fr", ".asso.fr", ".avocat.fr", ".chambagri.fr", ".chirurgiens-dentistes.fr", ".com.fr", ".experts-comptables.fr", ".geometre-expert.fr", ".gouv.fr", ".medecin.fr", ".nom.fr", ".notaires.fr", ".pharmacien.fr", ".port.fr", ".prd.fr", ".presse.fr", ".tm.fr", ".veterinaire.fr", ".com.ge", ".edu.ge", ".gov.ge", ".mil.ge", ".net.ge", ".org.ge", ".pvt.ge", ".co.gg", ".net.gg", ".org.gg", ".com.gh", ".edu.gh", ".gov.gh", ".org.gh", ".com.gi", ".gov.gi", ".ltd.gi", ".org.gi", ".co.gl", ".com.gl", ".edu.gl", ".net.gl", ".org.gl", ".com.gn", ".gov.gn", ".net.gn", ".org.gn", ".com.gp", ".mobi.gp", ".net.gp", ".org.gp", ".com.gr", ".edu.gr", ".net.gr", ".org.gr", ".com.gt", ".ind.gt", ".net.gt", ".org.gt", ".com.gu", ".co.gy", ".com.gy", ".net.gy", ".com.hk", ".edu.hk", ".gov.hk", ".idv.hk", ".inc.hk", ".ltd.hk", ".net.hk", ".org.hk", ".公司.hk", ".com.hn", ".edu.hn", ".net.hn", ".org.hn", ".com.hr", ".adult.ht", ".art.ht", ".asso.ht", ".com.ht", ".edu.ht", ".firm.ht", ".info.ht", ".net.ht", ".org.ht", ".perso.ht", ".pol.ht", ".pro.ht", ".rel.ht", ".shop.ht", ".2000.hu", ".agrar.hu", ".bolt.hu", ".casino.hu", ".city.hu", ".co.hu", ".erotica.hu", ".erotika.hu", ".film.hu", ".forum.hu", ".games.hu", ".hotel.hu", ".info.hu", ".ingatlan.hu", ".jogasz.hu", ".konyvelo.hu", ".lakas.hu", ".media.hu", ".news.hu", ".org.hu", ".priv.hu", ".reklam.hu", ".sex.hu", ".shop.hu", ".sport.hu", ".suli.hu", ".szex.hu", ".tm.hu", ".tozsde.hu", ".utazas.hu", ".video.hu", ".biz.id", ".co.id", ".my.id", ".or.id", ".web.id", ".ac.il", ".co.il", ".muni.il", ".net.il", ".org.il", ".ac.im", ".co.im", ".com.im", ".net.im", ".org.im", ".5g.in", ".6g.in", ".ahmdabad.in", ".ai.in", ".am.in", ".bihar.in", ".biz.in", ".business.in", ".ca.in", ".cn.in", ".co.in", ".com.in", ".coop.in", ".cs.in", ".delhi.in", ".dr.in", ".er.in", ".firm.in", ".gen.in", ".gujarat.in", ".ind.in", ".info.in", ".int.in", ".internet.in", ".io.in", ".me.in", ".net.in", ".org.in", ".pg.in", ".post.in", ".pro.in", ".travel.in", ".tv.in", ".uk.in", ".up.in", ".us.in", ".auz.info", ".com.iq", ".co.ir", ".abr.it", ".abruzzo.it", ".ag.it", ".agrigento.it", ".al.it", ".alessandria.it", ".alto-adige.it", ".altoadige.it", ".an.it", ".ancona.it", ".andria-barletta-trani.it", ".andria-trani-barletta.it", ".andriabarlettatrani.it", ".andriatranibarletta.it", ".ao.it", ".aosta.it", ".aoste.it", ".ap.it", ".aq.it", ".aquila.it", ".ar.it", ".arezzo.it", ".ascoli-piceno.it", ".ascolipiceno.it", ".asti.it", ".at.it", ".av.it", ".avellino.it", ".ba.it", ".balsan.it", ".bari.it", ".barletta-trani-andria.it", ".barlettatraniandria.it", ".bas.it", ".basilicata.it", ".belluno.it", ".benevento.it", ".bergamo.it", ".bg.it", ".bi.it", ".biella.it", ".bl.it", ".bn.it", ".bo.it", ".bologna.it", ".bolzano.it", ".bozen.it", ".br.it", ".brescia.it", ".brindisi.it", ".bs.it", ".bt.it", ".bz.it", ".ca.it", ".cagliari.it", ".cal.it", ".calabria.it", ".caltanissetta.it", ".cam.it", ".campania.it", ".campidano-medio.it", ".campidanomedio.it", ".campobasso.it", ".carbonia-iglesias.it", ".carboniaiglesias.it", ".carrara-massa.it", ".carraramassa.it", ".caserta.it", ".catania.it", ".catanzaro.it", ".cb.it", ".ce.it", ".cesena-forli.it", ".cesenaforli.it", ".ch.it", ".chieti.it", ".ci.it", ".cl.it", ".cn.it", ".co.it", ".como.it", ".cosenza.it", ".cr.it", ".cremona.it", ".crotone.it", ".cs.it", ".ct.it", ".cuneo.it", ".cz.it", ".dell-ogliastra.it", ".dellogliastra.it", ".emilia-romagna.it", ".emiliaromagna.it", ".emr.it", ".en.it", ".enna.it", ".fc.it", ".fe.it", ".fermo.it", ".ferrara.it", ".fg.it", ".fi.it", ".firenze.it", ".florence.it", ".fm.it", ".foggia.it", ".forli-cesena.it", ".forlicesena.it", ".fr.it", ".friuli-v-giulia.it", ".friuli-ve-giulia.it", ".friuli-vegiulia.it", ".friuli-venezia-giulia.it", ".friuli-veneziagiulia.it", ".friuli-vgiulia.it", ".friuliv-giulia.it", ".friulive-giulia.it", ".friulivegiulia.it", ".friulivenezia-giulia.it", ".friuliveneziagiulia.it", ".friulivgiulia.it", ".frosinone.it", ".fvg.it", ".ge.it", ".genoa.it", ".genova.it", ".go.it", ".gorizia.it", ".gr.it", ".grosseto.it", ".iglesias-carbonia.it", ".iglesiascarbonia.it", ".im.it", ".imperia.it", ".is.it", ".isernia.it", ".kr.it", ".la-spezia.it", ".laquila.it", ".laspezia.it", ".latina.it", ".laz.it", ".lazio.it", ".lc.it", ".le.it", ".lecce.it", ".lecco.it", ".li.it", ".lig.it", ".liguria.it", ".livorno.it", ".lo.it", ".lodi.it", ".lom.it", ".lombardia.it", ".lombardy.it", ".lt.it", ".lu.it", ".lucania.it", ".lucca.it", ".macerata.it", ".mantova.it", ".mar.it", ".marche.it", ".massa-carrara.it", ".massacarrara.it", ".matera.it", ".mb.it", ".mc.it", ".me.it", ".medio-campidano.it", ".mediocampidano.it", ".messina.it", ".mi.it", ".milan.it", ".milano.it", ".mn.it", ".mo.it", ".modena.it", ".mol.it", ".molise.it", ".monza-brianza.it", ".monza-e-della-brianza.it", ".monza.it", ".monzabrianza.it", ".monzaebrianza.it", ".monzaedellabrianza.it", ".ms.it", ".mt.it", ".na.it", ".naples.it", ".napoli.it", ".no.it", ".novara.it", ".nu.it", ".nuoro.it", ".og.it", ".ogliastra.it", ".olbia-tempio.it", ".olbiatempio.it", ".or.it", ".oristano.it", ".ot.it", ".pa.it", ".padova.it", ".padua.it", ".palermo.it", ".parma.it", ".pavia.it", ".pc.it", ".pd.it", ".pe.it", ".perugia.it", ".pesaro-urbino.it", ".pesarourbino.it", ".pescara.it", ".pg.it", ".pi.it", ".piacenza.it", ".piedmont.it", ".piemonte.it", ".pisa.it", ".pistoia.it", ".pmn.it", ".pn.it", ".po.it", ".pordenone.it", ".potenza.it", ".pr.it", ".prato.it", ".pt.it", ".pu.it", ".pug.it", ".puglia.it", ".pv.it", ".pz.it", ".ra.it", ".ragusa.it", ".ravenna.it", ".rc.it", ".re.it", ".reggio-calabria.it", ".reggio-emilia.it", ".reggiocalabria.it", ".reggioemilia.it", ".rg.it", ".ri.it", ".rieti.it", ".rimini.it", ".rm.it", ".rn.it", ".ro.it", ".roma.it", ".rome.it", ".rovigo.it", ".sa.it", ".salerno.it", ".sar.it", ".sardegna.it", ".sardinia.it", ".sassari.it", ".savona.it", ".si.it", ".sic.it", ".sicilia.it", ".sicily.it", ".siena.it", ".siracusa.it", ".so.it", ".sondrio.it", ".sp.it", ".sr.it", ".ss.it", ".suedtirol.it", ".sv.it", ".ta.it", ".taa.it", ".taranto.it", ".te.it", ".tempio-olbia.it", ".tempioolbia.it", ".teramo.it", ".terni.it", ".tn.it", ".to.it", ".torino.it", ".tos.it", ".toscana.it", ".tp.it", ".tr.it", ".trani-andria-barletta.it", ".trani-barletta-andria.it", ".traniandriabarletta.it", ".tranibarlettaandria.it", ".trapani.it", ".trentino-a-adige.it", ".trentino-aadige.it", ".trentino-alto-adige.it", ".trentino-altoadige.it", ".trentino-s-tirol.it", ".trentino-stirol.it", ".trentino-sud-tirol.it", ".trentino-sudtirol.it", ".trentino-sued-tirol.it", ".trentino-suedtirol.it", ".trentino.it", ".trentinoa-adige.it", ".trentinoaadige.it", ".trentinoalto-adige.it", ".trentinoaltoadige.it", ".trentinos-tirol.it", ".trentinosud-tirol.it", ".trentinosudtirol.it", ".trentinosued-tirol.it", ".trentinosuedtirol.it", ".trento.it", ".treviso.it", ".trieste.it", ".ts.it", ".turin.it", ".tuscany.it", ".tv.it", ".ud.it", ".udine.it", ".umb.it", ".umbria.it", ".urbino-pesaro.it", ".urbinopesaro.it", ".va.it", ".val-d-aosta.it", ".val-daosta.it", ".vald-aosta.it", ".valdaosta.it", ".valle-d-aosta.it", ".valle-daosta.it", ".valled-aosta.it", ".valledaosta.it", ".vao.it", ".varese.it", ".vb.it", ".vc.it", ".vda.it", ".ve.it", ".ven.it", ".veneto.it", ".venezia.it", ".venice.it", ".verbania.it", ".vercelli.it", ".verona.it", ".vi.it", ".vibo-valentia.it", ".vibovalentia.it", ".vicenza.it", ".viterbo.it", ".vr.it", ".vs.it", ".vt.it", ".vv.it", ".co.je", ".net.je", ".org.je", ".com.jm", ".net.jm", ".org.jm", ".com.jo", ".name.jo", ".net.jo", ".org.jo", ".sch.jo", ".akita.jp", ".co.jp", ".gr.jp", ".kyoto.jp", ".ne.jp", ".or.jp", ".osaka.jp", ".saga.jp", ".tokyo.jp", ".ac.ke", ".co.ke", ".go.ke", ".info.ke", ".me.ke", ".mobi.ke", ".ne.ke", ".or.ke", ".sc.ke", ".com.kg", ".net.kg", ".org.kg", ".com.kh", ".edu.kh", ".net.kh", ".org.kh", ".biz.ki", ".com.ki", ".edu.ki", ".gov.ki", ".info.ki", ".mobi.ki", ".net.ki", ".org.ki", ".phone.ki", ".tel.ki", ".com.km", ".nom.km", ".org.km", ".tm.km", ".com.kn", ".co.kr", ".go.kr", ".ms.kr", ".ne.kr", ".or.kr", ".pe.kr", ".re.kr", ".seoul.kr", ".com.kw", ".edu.kw", ".net.kw", ".org.kw", ".com.ky", ".net.ky", ".org.ky", ".com.kz", ".org.kz", ".com.lb", ".edu.lb", ".net.lb", ".org.lb", ".co.lc", ".com.lc", ".l.lc", ".net.lc", ".org.lc", ".p.lc", ".assn.lk", ".com.lk", ".edu.lk", ".grp.lk", ".hotel.lk", ".ltd.lk", ".ngo.lk", ".org.lk", ".soc.lk", ".web.lk", ".com.lr", ".org.lr", ".co.ls", ".net.ls", ".org.ls", ".asn.lv", ".com.lv", ".conf.lv", ".edu.lv", ".id.lv", ".mil.lv", ".net.lv", ".org.lv", ".com.ly", ".id.ly", ".med.ly", ".net.ly", ".org.ly", ".plc.ly", ".sch.ly", ".ac.ma", ".co.ma", ".net.ma", ".org.ma", ".press.ma", ".asso.mc", ".tm.mc", ".co.mg", ".com.mg", ".mil.mg", ".net.mg", ".nom.mg", ".org.mg", ".prd.mg", ".tm.mg", ".com.mk", ".edu.mk", ".inf.mk", ".name.mk", ".net.mk", ".org.mk", ".com.ml", ".biz.mm", ".com.mm", ".net.mm", ".org.mm", ".per.mm", ".com.mo", ".net.mo", ".org.mo", ".edu.mr", ".org.mr", ".perso.mr", ".co.ms", ".com.ms", ".org.ms", ".com.mt", ".net.mt", ".org.mt", ".ac.mu", ".co.mu", ".com.mu", ".net.mu", ".nom.mu", ".or.mu", ".org.mu", ".com.mv", ".ac.mw", ".co.mw", ".com.mw", ".coop.mw", ".edu.mw", ".int.mw", ".net.mw", ".org.mw", ".com.mx", ".edu.mx", ".gob.mx", ".net.mx", ".org.mx", ".com.my", ".mil.my", ".name.my", ".net.my", ".org.my", ".co.mz", ".edu.mz", ".net.mz", ".org.mz", ".alt.na", ".cc.na", ".co.na", ".com.na", ".edu.na", ".info.na", ".net.na", ".org.na", ".school.na", ".com.ne", ".info.ne", ".int.ne", ".org.ne", ".perso.ne", ".auz.net", ".gb.net", ".hu.net", ".in.net", ".jp.net", ".ru.net", ".se.net", ".uk.net", ".arts.nf", ".com.nf", ".firm.nf", ".info.nf", ".net.nf", ".org.nf", ".other.nf", ".per.nf", ".rec.nf", ".store.nf", ".web.nf", ".com.ng", ".edu.ng", ".gov.ng", ".i.ng", ".mobi.ng", ".name.ng", ".net.ng", ".org.ng", ".sch.ng", ".ac.ni", ".biz.ni", ".co.ni", ".com.ni", ".edu.ni", ".gob.ni", ".in.ni", ".info.ni", ".int.ni", ".mil.ni", ".net.ni", ".nom.ni", ".org.ni", ".web.ni", ".co.nl", ".com.nl", ".net.nl", ".co.no", ".fhs.no", ".folkebibl.no", ".fylkesbibl.no", ".gs.no", ".idrett.no", ".museum.no", ".priv.no", ".uenorge.no", ".vgs.no", ".aero.np", ".asia.np", ".biz.np", ".com.np", ".coop.np", ".info.np", ".mil.np", ".mobi.np", ".museum.np", ".name.np", ".net.np", ".org.np", ".pro.np", ".travel.np", ".biz.nr", ".com.nr", ".info.nr", ".net.nr", ".org.nr", ".co.nu", ".ac.nz", ".co.nz", ".geek.nz", ".gen.nz", ".iwi.nz", ".kiwi.nz", ".maori.nz", ".net.nz", ".org.nz", ".school.nz", ".biz.om", ".co.om", ".com.om", ".edu.om", ".gov.om", ".med.om", ".mil.om", ".museum.om", ".net.om", ".org.om", ".pro.om", ".sch.om", ".ae.org", ".hk.org", ".us.org", ".abo.pa", ".com.pa", ".edu.pa", ".gob.pa", ".ing.pa", ".med.pa", ".net.pa", ".nom.pa", ".org.pa", ".sld.pa", ".com.pe", ".edu.pe", ".gob.pe", ".mil.pe", ".net.pe", ".nom.pe", ".org.pe", ".asso.pf", ".com.pf", ".org.pf", ".com.pg", ".net.pg", ".org.pg", ".com.ph", ".net.ph", ".org.ph", ".biz.pk", ".com.pk", ".net.pk", ".org.pk", ".web.pk", ".agro.pl", ".aid.pl", ".atm.pl", ".augustow.pl", ".auto.pl", ".babia-gora.pl", ".bedzin.pl", ".beskidy.pl", ".bialowieza.pl", ".bialystok.pl", ".bielawa.pl", ".bieszczady.pl", ".biz.pl", ".boleslawiec.pl", ".bydgoszcz.pl", ".bytom.pl", ".cieszyn.pl", ".com.pl", ".czeladz.pl", ".czest.pl", ".dlugoleka.pl", ".edu.pl", ".elblag.pl", ".elk.pl", ".glogow.pl", ".gmina.pl", ".gniezno.pl", ".gorlice.pl", ".grajewo.pl", ".gsm.pl", ".ilawa.pl", ".info.pl", ".jaworzno.pl", ".jelenia-gora.pl", ".jgora.pl", ".kalisz.pl", ".karpacz.pl", ".kartuzy.pl", ".kaszuby.pl", ".katowice.pl", ".kazimierz-dolny.pl", ".kepno.pl", ".ketrzyn.pl", ".klodzko.pl", ".kobierzyce.pl", ".kolobrzeg.pl", ".konin.pl", ".konskowola.pl", ".kutno.pl", ".lapy.pl", ".lebork.pl", ".legnica.pl", ".lezajsk.pl", ".limanowa.pl", ".lomza.pl", ".lowicz.pl", ".lubin.pl", ".lukow.pl", ".mail.pl", ".malbork.pl", ".malopolska.pl", ".mazowsze.pl", ".mazury.pl", ".media.pl", ".miasta.pl", ".mielec.pl", ".mielno.pl", ".mil.pl", ".mragowo.pl", ".naklo.pl", ".net.pl", ".nieruchomosci.pl", ".nom.pl", ".nowaruda.pl", ".nysa.pl", ".olawa.pl", ".olecko.pl", ".olkusz.pl", ".olsztyn.pl", ".opoczno.pl", ".opole.pl", ".org.pl", ".ostroda.pl", ".ostroleka.pl", ".ostrowiec.pl", ".ostrowwlkp.pl", ".pc.pl", ".pila.pl", ".pisz.pl", ".podhale.pl", ".podlasie.pl", ".polkowice.pl", ".pomorskie.pl", ".pomorze.pl", ".powiat.pl", ".priv.pl", ".prochowice.pl", ".pruszkow.pl", ".przeworsk.pl", ".pulawy.pl", ".radom.pl", ".rawa-maz.pl", ".realestate.pl", ".rel.pl", ".rybnik.pl", ".rzeszow.pl", ".sanok.pl", ".sejny.pl", ".sex.pl", ".shop.pl", ".sklep.pl", ".skoczow.pl", ".slask.pl", ".slupsk.pl", ".sos.pl", ".sosnowiec.pl", ".stalowa-wola.pl", ".starachowice.pl", ".stargard.pl", ".suwalki.pl", ".swidnica.pl", ".swiebodzin.pl", ".swinoujscie.pl", ".szczecin.pl", ".szczytno.pl", ".szkola.pl", ".targi.pl", ".tarnobrzeg.pl", ".tgory.pl", ".tm.pl", ".tourism.pl", ".travel.pl", ".turek.pl", ".turystyka.pl", ".tychy.pl", ".ustka.pl", ".walbrzych.pl", ".warmia.pl", ".warszawa.pl", ".waw.pl", ".wegrow.pl", ".wielun.pl", ".wlocl.pl", ".wloclawek.pl", ".wodzislaw.pl", ".wolomin.pl", ".wroclaw.pl", ".zachpomor.pl", ".zagan.pl", ".zarow.pl", ".zgora.pl", ".zgorzelec.pl", ".co.pn", ".net.pn", ".org.pn", ".at.pr", ".biz.pr", ".ch.pr", ".com.pr", ".de.pr", ".eu.pr", ".fr.pr", ".info.pr", ".isla.pr", ".it.pr", ".name.pr", ".net.pr", ".nl.pr", ".org.pr", ".pro.pr", ".uk.pr", ".aaa.pro", ".aca.pro", ".acct.pro", ".arc.pro", ".avocat.pro", ".bar.pro", ".bus.pro", ".chi.pro", ".chiro.pro", ".cpa.pro", ".den.pro", ".dent.pro", ".eng.pro", ".jur.pro", ".law.pro", ".med.pro", ".min.pro", ".nur.pro", ".nurse.pro", ".pharma.pro", ".prof.pro", ".prx.pro", ".recht.pro", ".rel.pro", ".teach.pro", ".vet.pro", ".com.ps", ".net.ps", ".org.ps", ".co.pt", ".com.pt", ".org.pt", ".com.py", ".coop.py", ".edu.py", ".net.py", ".org.py", ".com.qa", ".edu.qa", ".mil.qa", ".name.qa", ".net.qa", ".org.qa", ".sch.qa", ".com.re", ".arts.ro", ".co.ro", ".com.ro", ".firm.ro", ".info.ro", ".ne.ro", ".nom.ro", ".nt.ro", ".or.ro", ".org.ro", ".rec.ro", ".sa.ro", ".srl.ro", ".store.ro", ".tm.ro", ".www.ro", ".co.rs", ".edu.rs", ".in.rs", ".org.rs", ".adygeya.ru", ".bashkiria.ru", ".bir.ru", ".cbg.ru", ".com.ru", ".dagestan.ru", ".grozny.ru", ".kalmykia.ru", ".kustanai.ru", ".marine.ru", ".mordovia.ru", ".msk.ru", ".mytis.ru", ".nalchik.ru", ".net.ru", ".nov.ru", ".org.ru", ".pp.ru", ".pyatigorsk.ru", ".spb.ru", ".vladikavkaz.ru", ".vladimir.ru", ".ac.rw", ".co.rw", ".net.rw", ".org.rw", ".com.sa", ".edu.sa", ".med.sa", ".net.sa", ".org.sa", ".pub.sa", ".sch.sa", ".com.sb", ".net.sb", ".org.sb", ".com.sc", ".net.sc", ".org.sc", ".com.sd", ".info.sd", ".net.sd", ".com.se", ".com.sg", ".edu.sg", ".net.sg", ".org.sg", ".per.sg", ".ae.si", ".at.si", ".cn.si", ".co.si", ".de.si", ".uk.si", ".us.si", ".com.sl", ".edu.sl", ".net.sl", ".org.sl", ".art.sn", ".com.sn", ".edu.sn", ".org.sn", ".perso.sn", ".univ.sn", ".com.so", ".net.so", ".org.so", ".biz.ss", ".com.ss", ".me.ss", ".net.ss", ".abkhazia.su", ".adygeya.su", ".aktyubinsk.su", ".arkhangelsk.su", ".armenia.su", ".ashgabad.su", ".azerbaijan.su", ".balashov.su", ".bashkiria.su", ".bryansk.su", ".bukhara.su", ".chimkent.su", ".dagestan.su", ".east-kazakhstan.su", ".exnet.su", ".georgia.su", ".grozny.su", ".ivanovo.su", ".jambyl.su", ".kalmykia.su", ".kaluga.su", ".karacol.su", ".karaganda.su", ".karelia.su", ".khakassia.su", ".krasnodar.su", ".kurgan.su", ".kustanai.su", ".lenug.su", ".mangyshlak.su", ".mordovia.su", ".msk.su", ".murmansk.su", ".nalchik.su", ".navoi.su", ".north-kazakhstan.su", ".nov.su", ".obninsk.su", ".penza.su", ".pokrovsk.su", ".sochi.su", ".spb.su", ".tashkent.su", ".termez.su", ".togliatti.su", ".troitsk.su", ".tselinograd.su", ".tula.su", ".tuva.su", ".vladikavkaz.su", ".vladimir.su", ".vologda.su", ".com.sv", ".edu.sv", ".gob.sv", ".org.sv", ".com.sy", ".co.sz", ".org.sz", ".com.tc", ".net.tc", ".org.tc", ".pro.tc", ".com.td", ".net.td", ".org.td", ".tourism.td", ".ac.th", ".co.th", ".in.th", ".or.th", ".ac.tj", ".aero.tj", ".biz.tj", ".co.tj", ".com.tj", ".coop.tj", ".dyn.tj", ".go.tj", ".info.tj", ".int.tj", ".mil.tj", ".museum.tj", ".my.tj", ".name.tj", ".net.tj", ".org.tj", ".per.tj", ".pro.tj", ".web.tj", ".com.tl", ".net.tl", ".org.tl", ".agrinet.tn", ".com.tn", ".defense.tn", ".edunet.tn", ".ens.tn", ".fin.tn", ".ind.tn", ".info.tn", ".intl.tn", ".nat.tn", ".net.tn", ".org.tn", ".perso.tn", ".rnrt.tn", ".rns.tn", ".rnu.tn", ".tourism.tn", ".av.tr", ".bbs.tr", ".biz.tr", ".com.tr", ".dr.tr", ".gen.tr", ".info.tr", ".name.tr", ".net.tr", ".org.tr", ".tel.tr", ".tv.tr", ".web.tr", ".biz.tt", ".co.tt", ".com.tt", ".info.tt", ".jobs.tt", ".mobi.tt", ".name.tt", ".net.tt", ".org.tt", ".pro.tt", ".club.tw", ".com.tw", ".ebiz.tw", ".game.tw", ".idv.tw", ".net.tw", ".org.tw", ".ac.tz", ".co.tz", ".go.tz", ".hotel.tz", ".info.tz", ".me.tz", ".mil.tz", ".mobi.tz", ".ne.tz", ".or.tz", ".sc.tz", ".tv.tz", ".biz.ua", ".cherkassy.ua", ".cherkasy.ua", ".chernigov.ua", ".chernivtsi.ua", ".chernovtsy.ua", ".ck.ua", ".cn.ua", ".co.ua", ".com.ua", ".crimea.ua", ".cv.ua", ".dn.ua", ".dnepropetrovsk.ua", ".dnipropetrovsk.ua", ".donetsk.ua", ".dp.ua", ".edu.ua", ".gov.ua", ".if.ua", ".in.ua", ".ivano-frankivsk.ua", ".kh.ua", ".kharkiv.ua", ".kharkov.ua", ".kherson.ua", ".khmelnitskiy.ua", ".kiev.ua", ".kirovograd.ua", ".km.ua", ".kr.ua", ".ks.ua", ".kyiv.ua", ".lg.ua", ".lt.ua", ".lugansk.ua", ".lutsk.ua", ".lviv.ua", ".mk.ua", ".net.ua", ".nikolaev.ua", ".od.ua", ".odesa.ua", ".odessa.ua", ".org.ua", ".pl.ua", ".poltava.ua", ".pp.ua", ".rivne.ua", ".rovno.ua", ".rv.ua", ".sebastopol.ua", ".sm.ua", ".sumy.ua", ".te.ua", ".ternopil.ua", ".uz.ua", ".uzhgorod.ua", ".vinnica.ua", ".vn.ua", ".volyn.ua", ".yalta.ua", ".zaporizhzhe.ua", ".zhitomir.ua", ".zp.ua", ".zt.ua", ".ac.ug", ".co.ug", ".com.ug", ".go.ug", ".ne.ug", ".or.ug", ".org.ug", ".sc.ug", ".ac.uk", ".co.uk", ".gov.uk", ".ltd.uk", ".me.uk", ".net.uk", ".org.uk", ".plc.uk", ".sch.uk", ".com.uy", ".edu.uy", ".net.uy", ".org.uy", ".biz.uz", ".co.uz", ".com.uz", ".net.uz", ".org.uz", ".com.vc", ".net.vc", ".org.vc", ".co.ve", ".com.ve", ".info.ve", ".net.ve", ".org.ve", ".web.ve", ".co.vi", ".com.vi", ".net.vi", ".org.vi", ".ac.vn", ".biz.vn", ".com.vn", ".edu.vn", ".gov.vn", ".health.vn", ".info.vn", ".int.vn", ".name.vn", ".net.vn", ".org.vn", ".pro.vn", ".com.vu", ".net.vu", ".org.vu", ".com.ws", ".net.ws", ".org.ws", ".com.ye", ".net.ye", ".org.ye", ".co.za", ".net.za", ".org.za", ".web.za", ".co.zm", ".com.zm", ".org.zm", ".co.zw", ".ком.рф", ".нет.рф", ".орг.рф", ".ак.срб", ".пр.срб", ".упр.срб", ".كمپنی.بھارت", ".कंपनी.भारत", ".কোম্পানি.ভারত", ".ਕੰਪਨੀ.ਭਾਰਤ", ".કંપની.ભારત", ".நிறுவனம்.இந்தியா", ".కంపెనీ.భారత్", ".ธุรกิจ.ไทย", ".個人.香港", ".公司.香港", ".政府.香港", ".教育.香港", ".組織.香港", ".網絡.香港", ".aaa", ".aarp", ".abarth", ".abb", ".abbott", ".abbvie", ".abc", ".able", ".abogado", ".abudhabi", ".ac", ".academy", ".accenture", ".accountant", ".accountants", ".aco", ".actor", ".ad", ".ads", ".adult", ".ae", ".aeg", ".aero", ".aetna", ".af", ".afl", ".africa", ".ag", ".agakhan", ".agency", ".ai", ".aig", ".airbus", ".airforce", ".airtel", ".akdn", ".al", ".alfaromeo", ".alibaba", ".alipay", ".allfinanz", ".allstate", ".ally", ".alsace", ".alstom", ".am", ".amazon", ".americanexpress", ".americanfamily", ".amex", ".amfam", ".amica", ".amsterdam", ".analytics", ".android", ".anquan", ".anz", ".ao", ".aol", ".apartments", ".app", ".apple", ".aq", ".aquarelle", ".ar", ".arab", ".aramco", ".archi", ".army", ".arpa", ".art", ".arte", ".as", ".asda", ".asia", ".associates", ".at", ".athleta", ".attorney", ".au", ".auction", ".audi", ".audible", ".audio", ".auspost", ".author", ".auto", ".autos", ".avianca", ".aw", ".aws", ".ax", ".axa", ".az", ".azure", ".ba", ".baby", ".baidu", ".banamex", ".bananarepublic", ".band", ".bank", ".bar", ".barcelona", ".barclaycard", ".barclays", ".barefoot", ".bargains", ".baseball", ".basketball", ".bauhaus", ".bayern", ".bb", ".bbc", ".bbt", ".bbva", ".bcg", ".bcn", ".bd", ".be", ".beats", ".beauty", ".beer", ".bentley", ".berlin", ".best", ".bestbuy", ".bet", ".bf", ".bg", ".bh", ".bharti", ".bi", ".bible", ".bid", ".bike", ".bing", ".bingo", ".bio", ".biz", ".bj", ".black", ".blackfriday", ".blockbuster", ".blog", ".bloomberg", ".blue", ".bm", ".bms", ".bmw", ".bn", ".bnpparibas", ".bo", ".boats", ".boehringer", ".bofa", ".bom", ".bond", ".boo", ".book", ".booking", ".bosch", ".bostik", ".boston", ".bot", ".boutique", ".box", ".br", ".bradesco", ".bridgestone", ".broadway", ".broker", ".brother", ".brussels", ".bs", ".bt", ".build", ".builders", ".business", ".buy", ".buzz", ".bv", ".bw", ".by", ".bz", ".bzh", ".ca", ".cab", ".cafe", ".cal", ".call", ".calvinklein", ".cam", ".camera", ".camp", ".canon", ".capetown", ".capital", ".capitalone", ".car", ".caravan", ".cards", ".care", ".career", ".careers", ".cars", ".casa", ".case", ".cash", ".casino", ".cat", ".catering", ".catholic", ".cba", ".cbn", ".cbre", ".cbs", ".cc", ".cd", ".center", ".ceo", ".cern", ".cf", ".cfa", ".cfd", ".cg", ".ch", ".chanel", ".channel", ".charity", ".chase", ".chat", ".cheap", ".chintai", ".christmas", ".chrome", ".church", ".ci", ".cipriani", ".circle", ".cisco", ".citadel", ".citi", ".citic", ".city", ".cityeats", ".ck", ".cl", ".claims", ".cleaning", ".click", ".clinic", ".clinique", ".clothing", ".cloud", ".club", ".clubmed", ".cm", ".cn", ".co", ".coach", ".codes", ".coffee", ".college", ".cologne", ".com", ".comcast", ".commbank", ".community", ".company", ".compare", ".computer", ".comsec", ".condos", ".construction", ".consulting", ".contact", ".contractors", ".cooking", ".cookingchannel", ".cool", ".coop", ".corsica", ".country", ".coupon", ".coupons", ".courses", ".cpa", ".cr", ".credit", ".creditcard", ".creditunion", ".cricket", ".crown", ".crs", ".cruise", ".cruises", ".cu", ".cuisinella", ".cv", ".cw", ".cx", ".cy", ".cymru", ".cyou", ".cz", ".dabur", ".dad", ".dance", ".data", ".date", ".dating", ".datsun", ".day", ".dclk", ".dds", ".de", ".deal", ".dealer", ".deals", ".degree", ".delivery", ".dell", ".deloitte", ".delta", ".democrat", ".dental", ".dentist", ".desi", ".design", ".dev", ".dhl", ".diamonds", ".diet", ".digital", ".direct", ".directory", ".discount", ".discover", ".dish", ".diy", ".dj", ".dk", ".dm", ".dnp", ".do", ".docs", ".doctor", ".dog", ".domains", ".dot", ".download", ".drive", ".dtv", ".dubai", ".dunlop", ".dupont", ".durban", ".dvag", ".dvr", ".dz", ".earth", ".eat", ".ec", ".eco", ".edeka", ".edu", ".education", ".ee", ".eg", ".email", ".emerck", ".energy", ".engineer", ".engineering", ".enterprises", ".epson", ".equipment", ".er", ".ericsson", ".erni", ".es", ".esq", ".estate", ".et", ".etisalat", ".eu", ".eurovision", ".eus", ".events", ".exchange", ".expert", ".exposed", ".express", ".extraspace", ".fage", ".fail", ".fairwinds", ".faith", ".family", ".fan", ".fans", ".farm", ".farmers", ".fashion", ".fast", ".fedex", ".feedback", ".ferrari", ".ferrero", ".fi", ".fiat", ".fidelity", ".fido", ".film", ".final", ".finance", ".financial", ".fire", ".firestone", ".firmdale", ".fish", ".fishing", ".fit", ".fitness", ".fj", ".fk", ".flickr", ".flights", ".flir", ".florist", ".flowers", ".fly", ".fm", ".fo", ".foo", ".food", ".foodnetwork", ".football", ".ford", ".forex", ".forsale", ".forum", ".foundation", ".fox", ".fr", ".free", ".fresenius", ".frl", ".frogans", ".frontdoor", ".frontier", ".ftr", ".fujitsu", ".fun", ".fund", ".furniture", ".futbol", ".fyi", ".ga", ".gal", ".gallery", ".gallo", ".gallup", ".game", ".games", ".gap", ".garden", ".gay", ".gb", ".gbiz", ".gd", ".gdn", ".ge", ".gea", ".gent", ".genting", ".george", ".gf", ".gg", ".ggee", ".gh", ".gi", ".gift", ".gifts", ".gives", ".giving", ".gl", ".glass", ".gle", ".global", ".globo", ".gm", ".gmail", ".gmbh", ".gmo", ".gmx", ".gn", ".godaddy", ".gold", ".goldpoint", ".golf", ".goo", ".goodyear", ".goog", ".google", ".gop", ".got", ".gov", ".gp", ".gq", ".gr", ".grainger", ".graphics", ".gratis", ".green", ".gripe", ".grocery", ".group", ".gs", ".gt", ".gu", ".guardian", ".gucci", ".guge", ".guide", ".guitars", ".guru", ".gw", ".gy", ".hair", ".hamburg", ".hangout", ".haus", ".hbo", ".hdfc", ".hdfcbank", ".health", ".healthcare", ".help", ".helsinki", ".here", ".hermes", ".hgtv", ".hiphop", ".hisamitsu", ".hitachi", ".hiv", ".hk", ".hkt", ".hm", ".hn", ".hockey", ".holdings", ".holiday", ".homedepot", ".homegoods", ".homes", ".homesense", ".honda", ".horse", ".hospital", ".host", ".hosting", ".hot", ".hoteles", ".hotels", ".hotmail", ".house", ".how", ".hr", ".hsbc", ".ht", ".hu", ".hughes", ".hyatt", ".hyundai", ".ibm", ".icbc", ".ice", ".icu", ".id", ".ie", ".ieee", ".ifm", ".ikano", ".il", ".im", ".imamat", ".imdb", ".immo", ".immobilien", ".in", ".inc", ".industries", ".infiniti", ".info", ".ing", ".ink", ".institute", ".insurance", ".insure", ".int", ".international", ".intuit", ".investments", ".io", ".ipiranga", ".iq", ".ir", ".irish", ".is", ".ismaili", ".ist", ".istanbul", ".it", ".itau", ".itv", ".jaguar", ".java", ".jcb", ".je", ".jeep", ".jetzt", ".jewelry", ".jio", ".jll", ".jm", ".jmp", ".jnj", ".jo", ".jobs", ".joburg", ".jot", ".joy", ".jp", ".jpmorgan", ".jprs", ".juegos", ".juniper", ".kaufen", ".kddi", ".ke", ".kerryhotels", ".kerrylogistics", ".kerryproperties", ".kfh", ".kg", ".kh", ".ki", ".kia", ".kids", ".kim", ".kinder", ".kindle", ".kitchen", ".kiwi", ".km", ".kn", ".koeln", ".komatsu", ".kosher", ".kp", ".kpmg", ".kpn", ".kr", ".krd", ".kred", ".kuokgroup", ".kw", ".ky", ".kyoto", ".kz", ".la", ".lacaixa", ".lamborghini", ".lamer", ".lancaster", ".lancia", ".land", ".landrover", ".lanxess", ".lasalle", ".lat", ".latino", ".latrobe", ".law", ".lawyer", ".lb", ".lc", ".lds", ".lease", ".leclerc", ".lefrak", ".legal", ".lego", ".lexus", ".lgbt", ".li", ".lidl", ".life", ".lifeinsurance", ".lifestyle", ".lighting", ".like", ".lilly", ".limited", ".limo", ".lincoln", ".linde", ".link", ".lipsy", ".live", ".living", ".lk", ".llc", ".llp", ".loan", ".loans", ".locker", ".locus", ".lol", ".london", ".lotte", ".lotto", ".love", ".lpl", ".lplfinancial", ".lr", ".ls", ".lt", ".ltd", ".ltda", ".lu", ".lundbeck", ".luxe", ".luxury", ".lv", ".ly", ".ma", ".madrid", ".maif", ".maison", ".makeup", ".man", ".management", ".mango", ".map", ".market", ".marketing", ".markets", ".marriott", ".marshalls", ".maserati", ".mattel", ".mba", ".mc", ".mckinsey", ".md", ".me", ".med", ".media", ".meet", ".melbourne", ".meme", ".memorial", ".men", ".menu", ".merckmsd", ".mg", ".mh", ".miami", ".microsoft", ".mil", ".mini", ".mint", ".mit", ".mitsubishi", ".mk", ".ml", ".mlb", ".mls", ".mm", ".mma", ".mn", ".mo", ".mobi", ".mobile", ".moda", ".moe", ".moi", ".mom", ".monash", ".money", ".monster", ".mormon", ".mortgage", ".moscow", ".moto", ".motorcycles", ".mov", ".movie", ".mp", ".mq", ".mr", ".ms", ".msd", ".mt", ".mtn", ".mtr", ".mu", ".museum", ".music", ".mutual", ".mv", ".mw", ".mx", ".my", ".mz", ".na", ".nab", ".nagoya", ".name", ".natura", ".navy", ".nba", ".nc", ".ne", ".nec", ".net", ".netbank", ".netflix", ".network", ".neustar", ".new", ".news", ".next", ".nextdirect", ".nexus", ".nf", ".nfl", ".ng", ".ngo", ".nhk", ".ni", ".nico", ".nike", ".nikon", ".ninja", ".nissan", ".nissay", ".nl", ".no", ".nokia", ".northwesternmutual", ".norton", ".now", ".nowruz", ".nowtv", ".np", ".nr", ".nra", ".nrw", ".ntt", ".nu", ".nyc", ".nz", ".obi", ".observer", ".office", ".okinawa", ".olayan", ".olayangroup", ".oldnavy", ".ollo", ".om", ".omega", ".one", ".ong", ".onl", ".online", ".ooo", ".open", ".oracle", ".orange", ".org", ".organic", ".origins", ".osaka", ".otsuka", ".ott", ".ovh", ".pa", ".page", ".panasonic", ".paris", ".pars", ".partners", ".parts", ".party", ".passagens", ".pay", ".pccw", ".pe", ".pet", ".pf", ".pfizer", ".pg", ".ph", ".pharmacy", ".phd", ".philips", ".phone", ".photo", ".photography", ".photos", ".physio", ".pics", ".pictet", ".pictures", ".pid", ".pin", ".ping", ".pink", ".pioneer", ".pizza", ".pk", ".pl", ".place", ".play", ".playstation", ".plumbing", ".plus", ".pm", ".pn", ".pnc", ".pohl", ".poker", ".politie", ".porn", ".post", ".pr", ".pramerica", ".praxi", ".press", ".prime", ".pro", ".prod", ".productions", ".prof", ".progressive", ".promo", ".properties", ".property", ".protection", ".pru", ".prudential", ".ps", ".pt", ".pub", ".pw", ".pwc", ".py", ".qa", ".qpon", ".quebec", ".quest", ".racing", ".radio", ".re", ".read", ".realestate", ".realtor", ".realty", ".recipes", ".red", ".redstone", ".redumbrella", ".rehab", ".reise", ".reisen", ".reit", ".reliance", ".ren", ".rent", ".rentals", ".repair", ".report", ".republican", ".rest", ".restaurant", ".review", ".reviews", ".rexroth", ".rich", ".richardli", ".ricoh", ".ril", ".rio", ".rip", ".ro", ".rocher", ".rocks", ".rodeo", ".rogers", ".room", ".rs", ".rsvp", ".ru", ".rugby", ".ruhr", ".run", ".rw", ".rwe", ".ryukyu", ".sa", ".saarland", ".safe", ".safety", ".sakura", ".sale", ".salon", ".samsclub", ".samsung", ".sandvik", ".sandvikcoromant", ".sanofi", ".sap", ".sarl", ".sas", ".save", ".saxo", ".sb", ".sbi", ".sbs", ".sc", ".sca", ".scb", ".schaeffler", ".schmidt", ".scholarships", ".school", ".schule", ".schwarz", ".science", ".scot", ".sd", ".se", ".search", ".seat", ".secure", ".security", ".seek", ".select", ".sener", ".services", ".seven", ".sew", ".sex", ".sexy", ".sfr", ".sg", ".sh", ".shangrila", ".sharp", ".shaw", ".shell", ".shia", ".shiksha", ".shoes", ".shop", ".shopping", ".shouji", ".show", ".showtime", ".si", ".silk", ".sina", ".singles", ".site", ".sj", ".sk", ".ski", ".skin", ".sky", ".skype", ".sl", ".sling", ".sm", ".smart", ".smile", ".sn", ".sncf", ".so", ".soccer", ".social", ".softbank", ".software", ".sohu", ".solar", ".solutions", ".song", ".sony", ".soy", ".spa", ".space", ".sport", ".spot", ".sr", ".srl", ".ss", ".st", ".stada", ".staples", ".star", ".statebank", ".statefarm", ".stc", ".stcgroup", ".stockholm", ".storage", ".store", ".stream", ".studio", ".study", ".style", ".su", ".sucks", ".supplies", ".supply", ".support", ".surf", ".surgery", ".suzuki", ".sv", ".swatch", ".swiss", ".sx", ".sy", ".sydney", ".systems", ".sz", ".tab", ".taipei", ".talk", ".taobao", ".target", ".tatamotors", ".tatar", ".tattoo", ".tax", ".taxi", ".tc", ".tci", ".td", ".tdk", ".team", ".tech", ".technology", ".tel", ".temasek", ".tennis", ".teva", ".tf", ".tg", ".th", ".thd", ".theater", ".theatre", ".tiaa", ".tickets", ".tienda", ".tiffany", ".tips", ".tires", ".tirol", ".tj", ".tjmaxx", ".tjx", ".tk", ".tkmaxx", ".tl", ".tm", ".tmall", ".tn", ".to", ".today", ".tokyo", ".tools", ".top", ".toray", ".toshiba", ".total", ".tours", ".town", ".toyota", ".toys", ".tr", ".trade", ".trading", ".training", ".travel", ".travelchannel", ".travelers", ".travelersinsurance", ".trust", ".trv", ".tt", ".tube", ".tui", ".tunes", ".tushu", ".tv", ".tvs", ".tw", ".tz", ".ua", ".ubank", ".ubs", ".ug", ".uk", ".unicom", ".university", ".uno", ".uol", ".ups", ".us", ".uy", ".uz", ".va", ".vacations", ".vana", ".vanguard", ".vc", ".ve", ".vegas", ".ventures", ".verisign", ".vermögensberater", ".vermögensberatung", ".versicherung", ".vet", ".vg", ".vi", ".viajes", ".video", ".vig", ".viking", ".villas", ".vin", ".vip", ".virgin", ".visa", ".vision", ".viva", ".vivo", ".vlaanderen", ".vn", ".vodka", ".volkswagen", ".volvo", ".vote", ".voting", ".voto", ".voyage", ".vu", ".vuelos", ".wales", ".walmart", ".walter", ".wang", ".wanggou", ".watch", ".watches", ".weather", ".weatherchannel", ".webcam", ".weber", ".website", ".wed", ".wedding", ".weibo", ".weir", ".wf", ".whoswho", ".wien", ".wiki", ".williamhill", ".win", ".windows", ".wine", ".winners", ".wme", ".wolterskluwer", ".woodside", ".work", ".works", ".world", ".wow", ".ws", ".wtc", ".wtf", ".xbox", ".xerox", ".xfinity", ".xihuan", ".xin", ".xxx", ".xyz", ".yachts", ".yahoo", ".yamaxun", ".yandex", ".ye", ".yodobashi", ".yoga", ".yokohama", ".you", ".youtube", ".yt", ".yun", ".za", ".zappos", ".zara", ".zero", ".zip", ".zm", ".zone", ".zuerich", ".zw", ".ελ", ".ευ", ".бг", ".бел", ".дети", ".ею", ".католик", ".ком", ".мкд", ".мон", ".москва", ".онлайн", ".орг", ".рус", ".рф", ".сайт", ".срб", ".укр", ".қаз", ".հայ", ".ישראל", ".קום", ".ابوظبي", ".اتصالات", ".ارامكو", ".الاردن", ".البحرين", ".الجزائر", ".السعودية", ".العليان", ".المغرب", ".امارات", ".ایران", ".بارت", ".بازار", ".بيتك", ".بھارت", ".تونس", ".سودان", ".سورية", ".شبكة", ".عراق", ".عرب", ".عمان", ".فلسطين", ".قطر", ".كاثوليك", ".كوم", ".مصر", ".مليسيا", ".موريتانيا", ".موقع", ".همراه", ".پاكستان", ".پاکستان", ".ڀارت", ".कॉम", ".नेट", ".भारत", ".भारतम्", ".भारोत", ".संगठन", ".বাংলা", ".ভারত", ".ভাৰত", ".ਭਾਰਤ", ".ભારત", ".ଭାରତ", ".இந்தியா", ".இலங்கை", ".சிங்கப்பூர்", ".భారత్", ".ಭಾರತ", ".ഭാരതം", ".ලංකා", ".คอม", ".ไทย", ".ລາວ", ".გე", ".みんな", ".アマゾン", ".クラウド", ".グーグル", ".コム", ".ストア", ".セール", ".ファッション", ".ポイント", ".世界", ".中信", ".中国", ".中國", ".中文网", ".亚马逊", ".企业", ".佛山", ".信息", ".健康", ".八卦", ".公司", ".公益", ".台湾", ".台灣", ".商城", ".商店", ".商标", ".嘉里", ".嘉里大酒店", ".在线", ".大拿", ".天主教", ".娱乐", ".家電", ".广东", ".微博", ".慈善", ".我爱你", ".手机", ".招聘", ".政务", ".政府", ".新加坡", ".新闻", ".时尚", ".書籍", ".机构", ".淡马锡", ".游戏", ".澳門", ".点看", ".移动", ".组织机构", ".网址", ".网店", ".网站", ".网络", ".联通", ".谷歌", ".购物", ".通販", ".集团", ".電訊盈科", ".飞利浦", ".食品", ".餐厅", ".香格里拉", ".香港", ".닷넷", ".닷컴", ".삼성", ".한국"}
	if len(strings.Split(domain, ".")) > 1 {
		for _, d := range domSuffix {
			if strings.HasSuffix(domain, d) {
				rootDomain, err := tld.Parse(fmt.Sprintf("https://%s", domain))
				if err != nil {
					return ""
				}
				return fmt.Sprintf("%s%s", rootDomain.Domain, d)
			}
		}
	}
	return ""
}

// About Site

// ParseUrl 格式化URL
func ParseUrl(uri string) (*UrlParse, error) {
	var res UrlParse
	var port string
	ps, err := url.Parse(uri)
	if err != nil && strings.HasPrefix(uri, "http") {
		tp1 := strings.Split(uri, "://")
		if len(tp1) == 2 {
			res.Scheme = tp1[0]
			hosts := tp1[1]
			if len(strings.Split(hosts, "/")) >= 2 {
				res.Hostname = strings.Split(hosts, "/")[0]
				if RootDomain(res.Hostname) == "" {
					return nil, err
				}
				if len(strings.Split(res.Hostname, ":")) >= 2 {
					res.Port = strings.Split(res.Hostname, ":")[0]

				} else {
					switch res.Scheme {
					case "http":
						res.Port = "80"
					case "https":
						res.Port = "443"
					}
					res.Port = "0"
				}
			} else {
				res.Hostname = hosts
			}
		} else {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	} else {
		res.Scheme = ps.Scheme
		res.Hostname = ps.Hostname()
		res.Path = ps.Path
		res.Query = ps.RawQuery
		port = ps.Port()
	}
	if port == "" {
		switch res.Scheme {
		case "https":
			res.Port = "443"
		case "ftp":
			res.Port = "21"
		default:
			res.Port = "80"
		}
	} else {
		res.Port = port
	}
	return &res, nil
}

// GetHeader 获取请求头
func GetHeader(args *GetHeaderArgs) http.Header {
	// 预设的User-Agent列表
	userAgents := []string{
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/126.0.0.0 Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/126.0.0.0 Safari/537.36 Edg/126.0.0.0",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/126.0.0.0 Safari/537.36 OPR/112.0.0.0",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.5735.289 Safari/537.36 ",
		"Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/94.0.4606.71 Safari/537.36 Core/1.94.202.400 QQBrowser/11.9.5355.400",
		"Mozilla/5.0 (Macintosh; U; Intel Mac OS X 10_6_8; en-us) AppleWebKit/534.50 (KHTML, like Gecko) Version/5.1 Safari/534.50",
		"Mozilla/5.0 (Windows; U; Windows NT 6.1; en-us) AppleWebKit/534.50 (KHTML, like Gecko) Version/5.1 Safari/534.50",
	}
	randIp := GetRandomIP()
	// 初始化HTTP头
	header := http.Header{}
	rand.Seed(time.Now().UnixNano())
	if args.Engine {
		header.Set("User-Agent", "Baiduspider+(+https://www.baidu.com/search/spider.htm);google|baiduspider|baidu|spider|sogou|bing|yahoo|soso|sosospider|360spider|youdao|jikeSpider;)")
	} else {
		header.Set("User-Agent", userAgents[rand.Intn(len(userAgents))])
	}
	header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	header.Set("Cache-Control", "no-cache")
	header.Set("Pragma", "no-cache")
	header.Set("CLIENT-IP", randIp)
	header.Set("X-Remote-Ip", randIp)
	header.Set("X-Remote-Addr", randIp)
	header.Set("X-Originating-Ip", randIp)
	header.Set("X-FORWARDED-FOR", randIp)
	header.Set("Connection", "keep-alive")
	header.Set("Upgrade-Insecure-Requests", "1")
	// 根据参数更新HTTP头
	if args.Switch == "Bing" {
		header.Set("Host", "cn.bing.com")
		header.Set("Referer", "https://www.bing.com")
	}
	if args.api != "" {
		if parse, err := url.Parse(args.api); err == nil {
			header.Set("Host", parse.Host) // 提取Host和Referer
		}
		header.Set("Referer", args.api)
	}
	// 传入的header参数
	for k, v := range args.header {
		header.Set(k, v)
	}
	return header
}

// HttpHeaderToMap http.header 转换成 map
func HttpHeaderToMap(header http.Header) map[string]string {
	var result = make(map[string]string)
	for k, vs := range header {
		for _, v := range vs {
			result[k] = v
		}
	}
	return result
}

// HttpHeaderToString http.header 转换成 String
func HttpHeaderToString(res *Response) string {
	Header := fmt.Sprintf("HTTP/%s.%s %s\n", strconv.Itoa(res.ProtoMinor), strconv.Itoa(res.ProtoMajor), res.Status)
	for k, vs := range res.Headers {
		for _, v := range vs {
			Header = fmt.Sprintf("%s%s: %s\n", Header, k, v)
		}
	}
	return Header
}

// CharSetContent 编码后的text
func CharSetContent(content []byte, body string, contentType string) (string, string) {
	var htmlEncode string
	// BOM检测编码
	if DetectBOM(content) {
		return Convert(body, "utf-8", "utf-8"), "utf-8"
	}
	// 响应头检测编码
	contentType = strings.ToLower(contentType)
	if contentType != "" {
		if strings.Contains(contentType, "gbk") || strings.Contains(contentType, "gb2312") || strings.Contains(contentType, "gb18030") || strings.Contains(contentType, "windows-1252") {
			htmlEncode = "gbk"
		} else if strings.Contains(contentType, "big5") {
			htmlEncode = "big5"
		} else if strings.Contains(contentType, "utf-8") {
			htmlEncode = "utf-8"
		}
	}
	if htmlEncode != "" {
		return Convert(body, htmlEncode, "utf-8"), htmlEncode
	}
	// 匹配正文中的编码
	match := regexp.MustCompile(`(?is)<meta[^>]*charset\s*=["']?\s*([A-Za-z0-9\-]+)`).FindStringSubmatch(body)
	if len(match) > 1 {
		contentType = strings.ToLower(match[1])
		if strings.Contains(contentType, "gbk") || strings.Contains(contentType, "gb2312") || strings.Contains(contentType, "gb18030") || strings.Contains(contentType, "windows-1252") {
			htmlEncode = "gbk"
		} else if strings.Contains(contentType, "big5") {
			htmlEncode = "big5"
		} else if strings.Contains(contentType, "utf-8") {
			htmlEncode = "utf-8"
		}
	}
	if htmlEncode != "" {
		return Convert(body, htmlEncode, "utf-8"), htmlEncode
	}
	// 自动检测编码
	_, contentType, _ = charset.DetermineEncoding(content, "")
	if contentType != "" {
		if strings.Contains(contentType, "gbk") || strings.Contains(contentType, "gb2312") || strings.Contains(contentType, "gb18030") || strings.Contains(contentType, "windows-1252") {
			htmlEncode = "gbk"
		} else if strings.Contains(contentType, "big5") {
			htmlEncode = "big5"
		} else if strings.Contains(contentType, "utf-8") {
			htmlEncode = "utf-8"
		}
	}
	if htmlEncode != "" {
		return Convert(body, htmlEncode, "utf-8"), htmlEncode
	}
	// 默认返回utf-8
	return Convert(body, "utf-8", "utf-8"), htmlEncode
}

// DetectBOM 检测是否为UTF-8 BOM
func DetectBOM(data []byte) bool {
	return bytes.HasPrefix(data, []byte{0xEF, 0xBB, 0xBF})
}

// Convert 编码HTML
func Convert(src string, srcCode string, tagCode string) string {
	if srcCode == tagCode {
		return src
	}
	srcCoder := mahonia.NewDecoder(srcCode)
	srcResult := srcCoder.ConvertString(src)
	tagCoder := mahonia.NewDecoder(tagCode)
	_, cdata, _ := tagCoder.Translate([]byte(srcResult), true)
	result := string(cdata)
	return result
}

// GetSiteBasic 网站提取标题
func GetSiteBasic(baseurl, text string) *SiteBasic {
	var result SiteBasic
	reTitle := regexp.MustCompile(`(?is)<title[^>]*>(.*?)</title>`)

	// 查找所有匹配项
	matchesTitle := reTitle.FindStringSubmatch(text)
	if len(matchesTitle) > 0 {
		title := matchesTitle[1]
		result.Title = title
	}

	reDes := regexp.MustCompile(`<meta\s+name=["']description["']\s+content=["'](.*?)["']\s*/?>`)
	// 查找所有匹配项
	matchesDes := reDes.FindStringSubmatch(text)
	if len(matchesDes) > 0 {
		desc := matchesDes[1]
		result.Description = desc
	}

	reKeywords := regexp.MustCompile(`<meta\s+name=["']keywords["']\s+content=["'](.*?)["']\s*/?>`)
	// 查找所有匹配项
	matchesKeywords := reKeywords.FindStringSubmatch(text)
	if len(matchesKeywords) > 0 {
		keywords := matchesKeywords[1]
		result.Keywords = keywords
	}

	favicon := GetFaviconPath(baseurl, text)
	result.Favicon = favicon
	return &result
}

// GetFaviconPath 获取favicon.ico的路径
func GetFaviconPath(uri, body string) string {
	//regFav := regexp.MustCompile(`href="(.*?favicon....)"`)
	regFav := regexp.MustCompile(`rel="icon" href="(.*?favicon[^"]*)">`)
	matchFav := regFav.FindAllStringSubmatch(body, -1)
	if len(matchFav) < 1 {
		regFav = regexp.MustCompile(`type="image/x-icon" href="(.*?favicon[^"]*)">`)
		matchFav = regFav.FindAllStringSubmatch(body, -1)
	}
	var faviconPath string
	u, err := url.Parse(uri)
	if err != nil {
		return ""
	}
	uri = fmt.Sprintf("%s://%s", u.Scheme, u.Host)
	if len(matchFav) > 0 {
		fav := matchFav[0][1]
		if fav[:2] == "//" {
			faviconPath = fmt.Sprintf("http:%s", fav)
		} else {
			if fav[:4] == "http" {
				faviconPath = fav
			} else {
				faviconPath = fmt.Sprintf("%s/%s", uri, strings.TrimPrefix(fav, "/"))
			}
		}
	} else {
		faviconPath = fmt.Sprintf("%s/favicon.ico", uri)
	}
	return faviconPath
}
