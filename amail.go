

/*2:*/


//line amail.w:31

//line license:1

// This file is part of Amail version 0.1
// Author Alexander Sychev
//
// Copyright (c) 2013 Alexander Sychev. All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are
// met:
//
//    * Redistributions of source code must retain the above copyright
// notice, this list of conditions and the following disclaimer.
//    * Redistributions in binary form must reproduce the above
// copyright notice, this list of conditions and the following disclaimer
// in the documentation and/or other materials provided with the
// distribution.
//    * The name of author may not be used to endorse or promote products derived from
// this software without specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
// "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
// LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
// A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
// OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
// SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
// LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
// DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
// THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

package main
//line amail.w:33
 import(


/*11:*/


//line amail.w:128

"flag"
"fmt"
"os"
"strings"
"sort"



/*:11*/



/*13:*/


//line amail.w:162

"unicode"
"unicode/utf8"



/*:13*/



/*15:*/


//line amail.w:180

"code.google.com/p/goplan9/plan9/client"
"github.com/golang/glog"



/*:15*/



/*28:*/


//line amail.w:308

"io"
"bufio"



/*:28*/



/*31:*/


//line amail.w:354

"bitbucket.org/santucco/goplumb"
"code.google.com/p/goplan9/plan9"



/*:31*/



/*38:*/


//line amail.w:505

"bitbucket.org/santucco/goacme"



/*:38*/



/*47:*/


//line amail.w:653

"strconv"



/*:47*/



/*168:*/


//line amail.w:2279

"os/exec"



/*:168*/



/*183:*/


//line amail.w:2604

"os/user"



/*:183*/



/*187:*/


//line amail.w:2625

"sync"



/*:187*/


//line amail.w:34

"time"
)

type(


/*19:*/


//line amail.w:204

mailbox struct{
name string


/*21:*/


//line amail.w:230

all messages
unread messages
mch chan int
dch chan int



/*:21*/



/*48:*/


//line amail.w:657

fid*client.Fid
total int



/*:48*/



/*74:*/


//line amail.w:950

shownew bool
showthreads bool
ech<-chan*goacme.Event
w*goacme.Window
cch chan bool



/*:74*/



/*79:*/


//line amail.w:1001

thread bool



/*:79*/



/*118:*/


//line amail.w:1530

rfch chan*struct{seek bool;msgs messages}
irfch chan*struct{seek bool;msgs messages}
reset bool



/*:118*/



/*128:*/


//line amail.w:1645

pos int



/*:128*/



/*140:*/


//line amail.w:1805

deleted int



/*:140*/



/*149:*/


//line amail.w:1982

lch chan[]int



/*:149*/


//line amail.w:207

}

mailboxes[]*mailbox

message struct{
id int


/*29:*/


//line amail.w:313

unread bool
box*mailbox



/*:29*/



/*66:*/


//line amail.w:883

deleted bool



/*:66*/



/*91:*/


//line amail.w:1210

from string
date time.Time
subject string



/*:91*/



/*100:*/


//line amail.w:1322

inreplyto string
messageid string



/*:100*/



/*103:*/


//line amail.w:1352

parent*message



/*:103*/



/*151:*/


//line amail.w:1990

w*goacme.Window



/*:151*/



/*159:*/


//line amail.w:2088

to[]string
cc[]string



/*:159*/



/*170:*/


//line amail.w:2301

text string
html string
showhtml bool
files[]*file
cids map[string]*file



/*:170*/


//line amail.w:214

}

messages[]*message




/*:19*/



/*98:*/


//line amail.w:1313

idmessages[]*message



/*:98*/



/*169:*/


//line amail.w:2283

file struct{
name string
mimetype string
path string
}



/*:169*/


//line amail.w:39

)



/*41:*/


//line amail.w:534

const mailboxfmt= "%-30s\t%10d\t%10d\n"



/*:41*/



/*45:*/


//line amail.w:613

const mailboxfmtprc= "%-30s\t%10d\t%10d\t%d%%\n"



/*:45*/


//line amail.w:42


var(


/*3:*/


//line amail.w:77

exit chan bool= make(chan bool)



/*:3*/



/*6:*/


//line amail.w:95

wch chan int= make(chan int,100)
wcount int



/*:6*/



/*10:*/


//line amail.w:120

shownew bool
showthreads bool
levelmark string
newmark string
skipboxes[]string



/*:10*/



/*16:*/


//line amail.w:185

fsys*client.Fsys
rfid*client.Fid
srv string= "mail"



/*:16*/



/*20:*/


//line amail.w:221

boxes mailboxes



/*:20*/



/*23:*/


//line amail.w:248

mch= make(chan*struct{name string;id int},100)
dch= make(chan*struct{name string;id int},100)
bch= make(chan string,10)
rfch= make(chan*mailbox,100)



/*:23*/



/*39:*/


//line amail.w:509

mw*goacme.Window
ech<-chan*goacme.Event



/*:39*/



/*44:*/


//line amail.w:609

shown= make(map[string]int)



/*:44*/



/*83:*/


//line amail.w:1143

lch= make(chan*map[string][]int,100)



/*:83*/



/*89:*/


//line amail.w:1194

deleted= "(deleted)-"



/*:89*/



/*99:*/


//line amail.w:1317

idmap= make(map[string]*struct{msg*message;children idmessages})
idch= make(chan struct{id string;val interface{}},100)



/*:99*/



/*175:*/


//line amail.w:2442

home string



/*:175*/



/*184:*/


//line amail.w:2608

cuser string



/*:184*/



/*188:*/


//line amail.w:2629

once sync.Once



/*:188*/



/*198:*/


//line amail.w:2914

plan9dir string



/*:198*/


//line amail.w:45

debug glog.Level= 1
)

func main(){
glog.V(debug).Infoln("main")
defer glog.V(debug).Infoln("main is done")


/*12:*/


//line amail.w:136

{
glog.V(debug).Infoln("parsing command line arguments")
var skip string
flag.BoolVar(&shownew,"new",true,"show new messages only")
flag.BoolVar(&showthreads,"threads",true,"show threads of messages")
flag.StringVar(&skip,"skip","","boxes to be skiped, separated by comma")
flag.StringVar(&levelmark,"levelmark","+","mark of level for threads")
flag.StringVar(&newmark,"newmark","(*)","mark of new messages")
flag.Usage= func(){
fmt.Fprintf(os.Stderr,"Mail client for Acme programming environment\n")
fmt.Fprintf(os.Stderr,"Usage: %s [options][<mailbox>]\n",os.Args[0])
fmt.Fprintln(os.Stderr,"Options:")
flag.PrintDefaults()
}
flag.Parse()


/*14:*/


//line amail.w:168

glog.V(debug).Infoln("checking of levelmark and newmark")
if r,_:=utf8.DecodeLastRuneInString(levelmark);unicode.IsDigit(r){
fmt.Fprintln(os.Stderr,"last symbol of level mark shouldn't be a digit")
os.Exit(1)
}
if r,_:=utf8.DecodeRuneInString(newmark);unicode.IsDigit(r){
fmt.Fprintln(os.Stderr,"first symbol of new mark shouldn't be a digit")
os.Exit(1)
}



/*:14*/


//line amail.w:152

if len(skip)> 0{
skipboxes= strings.Split(skip,", ")
sort.Strings(skipboxes)
glog.V(debug).Infof("these mailboxes will be skipped: %v\n",skipboxes)

}
}



/*:12*/


//line amail.w:52



/*17:*/


//line amail.w:191

{
glog.V(debug).Infoln("try to open mailfs")
var err error
if fsys,err= client.MountService(srv);err!=nil{
glog.Errorf("can't mount mailfs: %v\n",err)
os.Exit(1)
}
}




/*:17*/


//line amail.w:53



/*32:*/


//line amail.w:361

{
glog.V(debug).Infoln("trying to open 'seemail' plumbing port")
if sm,err:=goplumb.Open("seemail",plan9.OREAD);err!=nil{
glog.Errorf("can't open plumb/seemail: %s\n",err)
}else{
sch,err:=sm.MessageChannel(0)
if err!=nil{
glog.Errorf("can't get message channal for plumb/seemail: %s\n",err)
}else{
go func(){
defer sm.Close()
defer glog.V(debug).Infoln("plumbing goroutine is done")
for{
select{


/*4:*/


//line amail.w:81

case<-exit:
glog.V(debug).Infoln("on exit!")
return



/*:4*/


//line amail.w:376

case m,ok:=<-sch:
if!ok{
glog.Warningln("it seems plumber has finished")
sch= nil
continue
}
glog.V(debug).Infof("a plumbing message has been received: %v\n",m)
if m.Attr["filetype"]!="mail"{
glog.Warningln("attribute 'filetype' is not 'mail'")
continue
}
v,ok:=m.Attr["mailtype"]
if!ok{
glog.Warningln("can't find 'mailtype' attribute")
continue
}
b:=strings.Split(string(m.Data),"/")
if len(b)<3{
glog.Warningln("can't read a name of mailbox and a number of message")
continue
}
num,err:=strconv.Atoi(b[2])
if err!=nil{
glog.Error(err)
continue
}
if v=="new"{
glog.V(debug).Infof("'%d' is a new message in the '%s' mailbox\n",num,b[1])
mch<-&struct{name string;id int}{name:b[1],id:num}
}else if v=="delete"{
glog.V(debug).Infof("'%d' is a deleted message in the '%s' mailbox\n",num,b[1])
dch<-&struct{name string;id int}{name:b[1],id:num}
}
}
}
}()
}
}
}



/*:32*/


//line amail.w:54



/*25:*/


//line amail.w:272

glog.V(debug).Infoln("initialization of root of mailfs")
var err error
rfid,err= fsys.Walk(".")
if err!=nil{
glog.Errorf("can't open mailfs: %v\n",err)
os.Exit(1)
}
defer rfid.Close()




/*:25*/


//line amail.w:55

if len(flag.Args())> 0{


/*34:*/


//line amail.w:456

go func(){
glog.V(debug).Infoln("Start a main message loop")
defer glog.V(debug).Infoln("main message loop is done")
for{
select{


/*4:*/


//line amail.w:81

case<-exit:
glog.V(debug).Infoln("on exit!")
return



/*:4*/


//line amail.w:462

case b:=<-rfch:
if b==nil{


/*42:*/


//line amail.w:538

if mw!=nil{
glog.V(debug).Infoln("printing of the mailboxes")
if err:=mw.WriteAddr("0,$");err!=nil{
glog.Errorf("can't write to 'addr' file: %s\n",err)
}else if data,err:=mw.File("data");err!=nil{
glog.Errorf("can't open 'data' file: %s\n",err)
}else{
for _,v:=range boxes{
data.Write([]byte(fmt.Sprintf(mailboxfmt,v.name,len(v.unread),len(v.all))))
}
}
w:=mw


/*50:*/


//line amail.w:718

if w!=nil{
glog.V(debug).Infoln("setting the window to clean state")
if err:=w.WriteCtl("clean");err!=nil{
glog.Errorf("can't write to 'ctl' file: %s\n",err)
}
}



/*:50*/


//line amail.w:551

}




/*:42*/


//line amail.w:465

}else{


/*46:*/


//line amail.w:618

glog.V(debug).Infof("refreshing main window for the '%s' mailbox, len(all): %d, total: %d\n",b.name,len(b.all),b.total)
if mw!=nil{
if len(b.all)!=b.total{
if c,ok:=shown[b.name];!ok||c<99{
shown[b.name]= c+1
continue
}else{
shown[b.name]= 0
}
}

if err:=mw.WriteAddr("0/^%s.*\\n/",escape(b.name));err!=nil{
glog.Errorf("can't write to 'addr' file: %s\n",err)
continue
}

if data,err:=mw.File("data");err!=nil{
glog.Errorf("can't open 'data' file: %s\n",err)
}else if len(b.all)!=b.total{
if _,err:=data.Write([]byte(fmt.Sprintf(mailboxfmtprc,b.name,len(b.unread),len(b.all),len(b.all)*100/b.total)));
err!=nil{
glog.Errorf("can't write to 'data' file: %s\n",err)
continue
}
}else if _,err:=data.Write([]byte(fmt.Sprintf(mailboxfmt,b.name,len(b.unread),len(b.all))));
err!=nil{
glog.Errorf("can't write to 'data' file: %s\n",err)
continue
}
w:=mw


/*50:*/


//line amail.w:718

if w!=nil{
glog.V(debug).Infoln("setting the window to clean state")
if err:=w.WriteCtl("clean");err!=nil{
glog.Errorf("can't write to 'ctl' file: %s\n",err)
}
}



/*:50*/


//line amail.w:649

}



/*:46*/


//line amail.w:467

}


/*43:*/


//line amail.w:557

case ev,ok:=<-ech:
glog.V(debug).Infof("an event from main window has been received: %v\n",ev)
if!ok{
ech= nil


/*9:*/


//line amail.w:114

glog.V(debug).Infoln("decrease the windows count")
wch<-(-1)




/*:9*/


//line amail.w:562

return
}
if(ev.Type&goacme.Execute)==goacme.Execute{
switch ev.Text{
case"ShowNew":
shownew= true
case"ShowAll":
shownew= false
case"ShowThreads":
showthreads= true
case"ShowPlain":
showthreads= false
case"Del":
mw.UnreadEvent(ev)
mw.Close()
mw= nil


/*5:*/


//line amail.w:87

glog.V(debug).Infoln("exit!")
close(exit)



/*:5*/


//line amail.w:579

case"debug":
debug= 0
continue
case"nodebug":
debug= 1
continue
default:
mw.UnreadEvent(ev)
continue
}


/*138:*/


//line amail.w:1778

glog.V(debug).Infoln("writing a tag of the main window")
if err:=writeTag(mw,fmt.Sprintf(" %s %s ",
func()string{if shownew{return"ShowAll"}else{return"ShowNew"}}(),
func()string{if showthreads{return"ShowPlain"}else{return"ShowThreads"}}()));
err!=nil{
glog.Errorf("can't set a tag of the main window: %v",err)
}



/*:138*/


//line amail.w:590

continue
}else if(ev.Type&goacme.Look)==goacme.Look{
name:=ev.Text
// a box name can contain spaces
if len(ev.Arg)> 0{
name+= " "+ev.Arg
}
if i,ok:=boxes.Search(name);ok{
box:=boxes[i]


/*76:*/


//line amail.w:964

glog.V(debug).Infof("inform the '%s' mailbox to create a window\n",box.name)
box.cch<-true



/*:76*/


//line amail.w:600

continue
}
}
mw.UnreadEvent(ev)




/*:43*/


//line amail.w:469

}
}
}()



/*:34*/


//line amail.w:57

for _,name:=range flag.Args(){


/*35:*/


//line amail.w:475

glog.V(debug).Infof("creating a '%s' mailbox\n",name)
box:=&mailbox{name:name,

/*22:*/


//line amail.w:237

mch:make(chan int,100),
dch:make(chan int,100),




/*:22*/



/*75:*/


//line amail.w:958

shownew:shownew,
showthreads:showthreads,
cch:make(chan bool,100),



/*:75*/



/*119:*/


//line amail.w:1536

rfch:make(chan*struct{seek bool;msgs messages},100),
irfch:make(chan*struct{seek bool;msgs messages},100),



/*:119*/



/*150:*/


//line amail.w:1986

lch:make(chan[]int,100),



/*:150*/


//line amail.w:477
}
boxes= append(boxes,box)
sort.Sort(boxes)



/*:35*/


//line amail.w:59



/*63:*/


//line amail.w:837

go box.loop()



/*:63*/


//line amail.w:60



/*76:*/


//line amail.w:964

glog.V(debug).Infof("inform the '%s' mailbox to create a window\n",box.name)
box.cch<-true



/*:76*/


//line amail.w:61

}
}else{


/*40:*/


//line amail.w:514

glog.V(debug).Infoln("creating the main window")
defer goacme.DeleteAll()

var err error
if mw,err= goacme.New();err!=nil{
glog.Errorf("can't create a window: %v\n",err)
os.Exit(1)
}
name:="Amail"
w:=mw


/*53:*/


//line amail.w:748

glog.V(debug).Infoln("printing a name for a window")
if err:=w.WriteCtl("name %s",strings.Replace(name," ","␣",-1));err!=nil{
glog.Errorf("can't write to 'ctl' file: %s\n",err)
}



/*:53*/


//line amail.w:525

if ech,err= mw.EventChannel(0,goacme.Mouse,goacme.Look|goacme.Execute);err!=nil{
glog.Errorf("can't open an event channel of the window %v\n",err)
os.Exit(1)
}


/*138:*/


//line amail.w:1778

glog.V(debug).Infoln("writing a tag of the main window")
if err:=writeTag(mw,fmt.Sprintf(" %s %s ",
func()string{if shownew{return"ShowAll"}else{return"ShowNew"}}(),
func()string{if showthreads{return"ShowPlain"}else{return"ShowThreads"}}()));
err!=nil{
glog.Errorf("can't set a tag of the main window: %v",err)
}



/*:138*/


//line amail.w:530



/*8:*/


//line amail.w:109

glog.V(debug).Infoln("increase the windows count")
wch<-1



/*:8*/


//line amail.w:531




/*:40*/


//line amail.w:64



/*34:*/


//line amail.w:456

go func(){
glog.V(debug).Infoln("Start a main message loop")
defer glog.V(debug).Infoln("main message loop is done")
for{
select{


/*4:*/


//line amail.w:81

case<-exit:
glog.V(debug).Infoln("on exit!")
return



/*:4*/


//line amail.w:462

case b:=<-rfch:
if b==nil{


/*42:*/


//line amail.w:538

if mw!=nil{
glog.V(debug).Infoln("printing of the mailboxes")
if err:=mw.WriteAddr("0,$");err!=nil{
glog.Errorf("can't write to 'addr' file: %s\n",err)
}else if data,err:=mw.File("data");err!=nil{
glog.Errorf("can't open 'data' file: %s\n",err)
}else{
for _,v:=range boxes{
data.Write([]byte(fmt.Sprintf(mailboxfmt,v.name,len(v.unread),len(v.all))))
}
}
w:=mw


/*50:*/


//line amail.w:718

if w!=nil{
glog.V(debug).Infoln("setting the window to clean state")
if err:=w.WriteCtl("clean");err!=nil{
glog.Errorf("can't write to 'ctl' file: %s\n",err)
}
}



/*:50*/


//line amail.w:551

}




/*:42*/


//line amail.w:465

}else{


/*46:*/


//line amail.w:618

glog.V(debug).Infof("refreshing main window for the '%s' mailbox, len(all): %d, total: %d\n",b.name,len(b.all),b.total)
if mw!=nil{
if len(b.all)!=b.total{
if c,ok:=shown[b.name];!ok||c<99{
shown[b.name]= c+1
continue
}else{
shown[b.name]= 0
}
}

if err:=mw.WriteAddr("0/^%s.*\\n/",escape(b.name));err!=nil{
glog.Errorf("can't write to 'addr' file: %s\n",err)
continue
}

if data,err:=mw.File("data");err!=nil{
glog.Errorf("can't open 'data' file: %s\n",err)
}else if len(b.all)!=b.total{
if _,err:=data.Write([]byte(fmt.Sprintf(mailboxfmtprc,b.name,len(b.unread),len(b.all),len(b.all)*100/b.total)));
err!=nil{
glog.Errorf("can't write to 'data' file: %s\n",err)
continue
}
}else if _,err:=data.Write([]byte(fmt.Sprintf(mailboxfmt,b.name,len(b.unread),len(b.all))));
err!=nil{
glog.Errorf("can't write to 'data' file: %s\n",err)
continue
}
w:=mw


/*50:*/


//line amail.w:718

if w!=nil{
glog.V(debug).Infoln("setting the window to clean state")
if err:=w.WriteCtl("clean");err!=nil{
glog.Errorf("can't write to 'ctl' file: %s\n",err)
}
}



/*:50*/


//line amail.w:649

}



/*:46*/


//line amail.w:467

}


/*43:*/


//line amail.w:557

case ev,ok:=<-ech:
glog.V(debug).Infof("an event from main window has been received: %v\n",ev)
if!ok{
ech= nil


/*9:*/


//line amail.w:114

glog.V(debug).Infoln("decrease the windows count")
wch<-(-1)




/*:9*/


//line amail.w:562

return
}
if(ev.Type&goacme.Execute)==goacme.Execute{
switch ev.Text{
case"ShowNew":
shownew= true
case"ShowAll":
shownew= false
case"ShowThreads":
showthreads= true
case"ShowPlain":
showthreads= false
case"Del":
mw.UnreadEvent(ev)
mw.Close()
mw= nil


/*5:*/


//line amail.w:87

glog.V(debug).Infoln("exit!")
close(exit)



/*:5*/


//line amail.w:579

case"debug":
debug= 0
continue
case"nodebug":
debug= 1
continue
default:
mw.UnreadEvent(ev)
continue
}


/*138:*/


//line amail.w:1778

glog.V(debug).Infoln("writing a tag of the main window")
if err:=writeTag(mw,fmt.Sprintf(" %s %s ",
func()string{if shownew{return"ShowAll"}else{return"ShowNew"}}(),
func()string{if showthreads{return"ShowPlain"}else{return"ShowThreads"}}()));
err!=nil{
glog.Errorf("can't set a tag of the main window: %v",err)
}



/*:138*/


//line amail.w:590

continue
}else if(ev.Type&goacme.Look)==goacme.Look{
name:=ev.Text
// a box name can contain spaces
if len(ev.Arg)> 0{
name+= " "+ev.Arg
}
if i,ok:=boxes.Search(name);ok{
box:=boxes[i]


/*76:*/


//line amail.w:964

glog.V(debug).Infof("inform the '%s' mailbox to create a window\n",box.name)
box.cch<-true



/*:76*/


//line amail.w:600

continue
}
}
mw.UnreadEvent(ev)




/*:43*/


//line amail.w:469

}
}
}()



/*:34*/


//line amail.w:65

go func(){


/*26:*/


//line amail.w:284

{
glog.V(debug).Infoln("enumerating of mailboxes")
fi,err:=rfid.Dirreadall()
if err!=nil{
glog.Errorf("can't read mailfs: %v\n",err)


/*5:*/


//line amail.w:87

glog.V(debug).Infoln("exit!")
close(exit)



/*:5*/


//line amail.w:290

return
}
for _,f:=range fi{
if f.Mode&plan9.DMDIR==plan9.DMDIR{
name:=f.Name


/*27:*/


//line amail.w:303

glog.V(debug).Infof("send a mailbox '%s' to put in the list\n",name)
bch<-name



/*:27*/


//line amail.w:296

}
}
glog.V(debug).Infoln("enumerating of mailboxes is done")
}



/*:26*/


//line amail.w:67

}()
}


/*33:*/


//line amail.w:422

glog.V(debug).Infoln("process events are specific for the list of mailboxes")
for{
select{


/*4:*/


//line amail.w:81

case<-exit:
glog.V(debug).Infoln("on exit!")
return



/*:4*/


//line amail.w:426

case name,ok:=<-bch:
if!ok{
bch= nil


/*73:*/


//line amail.w:938

glog.V(debug).Infoln("sending a signal to refresh all mailboxes")
rfch<-nil



/*:73*/


//line amail.w:430

continue
}


/*54:*/


//line amail.w:755

glog.V(debug).Infoln("continue if the box should be skiped")
if i:=sort.SearchStrings(skipboxes,name);i!=len(skipboxes)&&skipboxes[i]==name{
continue
}



/*:54*/


//line amail.w:433



/*35:*/


//line amail.w:475

glog.V(debug).Infof("creating a '%s' mailbox\n",name)
box:=&mailbox{name:name,

/*22:*/


//line amail.w:237

mch:make(chan int,100),
dch:make(chan int,100),




/*:22*/



/*75:*/


//line amail.w:958

shownew:shownew,
showthreads:showthreads,
cch:make(chan bool,100),



/*:75*/



/*119:*/


//line amail.w:1536

rfch:make(chan*struct{seek bool;msgs messages},100),
irfch:make(chan*struct{seek bool;msgs messages},100),



/*:119*/



/*150:*/


//line amail.w:1986

lch:make(chan[]int,100),



/*:150*/


//line amail.w:477
}
boxes= append(boxes,box)
sort.Sort(boxes)



/*:35*/


//line amail.w:434



/*73:*/


//line amail.w:938

glog.V(debug).Infoln("sending a signal to refresh all mailboxes")
rfch<-nil



/*:73*/


//line amail.w:435



/*63:*/


//line amail.w:837

go box.loop()



/*:63*/


//line amail.w:436

case d:=<-mch:
name:=d.name


/*37:*/


//line amail.w:494

glog.V(debug).Infof("looking for a '%s' mailbox\n",name)
i,ok:=boxes.Search(name)
if!ok{


/*54:*/


//line amail.w:755

glog.V(debug).Infoln("continue if the box should be skiped")
if i:=sort.SearchStrings(skipboxes,name);i!=len(skipboxes)&&skipboxes[i]==name{
continue
}



/*:54*/


//line amail.w:498

glog.Warningf("can't find message box '%s'\n",name)
continue
}




/*:37*/


//line amail.w:439

glog.V(debug).Infof("sending '%d' to add in the '%s' mailbox\n",d.id,boxes[i].name)
boxes[i].mch<-d.id
case d:=<-dch:
name:=d.name


/*37:*/


//line amail.w:494

glog.V(debug).Infof("looking for a '%s' mailbox\n",name)
i,ok:=boxes.Search(name)
if!ok{


/*54:*/


//line amail.w:755

glog.V(debug).Infoln("continue if the box should be skiped")
if i:=sort.SearchStrings(skipboxes,name);i!=len(skipboxes)&&skipboxes[i]==name{
continue
}



/*:54*/


//line amail.w:498

glog.Warningf("can't find message box '%s'\n",name)
continue
}




/*:37*/


//line amail.w:444

glog.V(debug).Infof("sending '%d' to delete from the '%s' mailbox\n",d.id,boxes[i].name)
boxes[i].dch<-d.id


/*7:*/


//line amail.w:100

case i:=<-wch:
wcount+= i
if wcount==0{


/*5:*/


//line amail.w:87

glog.V(debug).Infoln("exit!")
close(exit)



/*:5*/


//line amail.w:104

return
}



/*:7*/



/*86:*/


//line amail.w:1156

case d:=<-lch:
if d==nil{
continue
}
for name,ids:=range*d{


/*37:*/


//line amail.w:494

glog.V(debug).Infof("looking for a '%s' mailbox\n",name)
i,ok:=boxes.Search(name)
if!ok{


/*54:*/


//line amail.w:755

glog.V(debug).Infoln("continue if the box should be skiped")
if i:=sort.SearchStrings(skipboxes,name);i!=len(skipboxes)&&skipboxes[i]==name{
continue
}



/*:54*/


//line amail.w:498

glog.Warningf("can't find message box '%s'\n",name)
continue
}




/*:37*/


//line amail.w:1162

boxes[i].lch<-ids
}



/*:86*/



/*102:*/


//line amail.w:1341

case v:=<-idch:
if v.val==nil{


/*112:*/


//line amail.w:1448

{
val,ok:=idmap[v.id]
if!ok{
continue
}
delete(idmap,v.id)
for i,v:=range val.children{
glog.V(debug).Infof("clear the parent of the '%d'\n",v.id)
val.children[i].parent= nil
}
if val.msg!=nil&&val.msg.parent!=nil{
if p,ok:=idmap[val.msg.parent.messageid];ok{
for i,_:=range p.children{
if p.children[i]==val.msg{
glog.V(debug).Infof("removing the '%d' message from the children of the message '%d'\n",
val.msg.id,val.msg.parent.id)
p.children.Delete(i)
break
}
}
}
}
}



/*:112*/


//line amail.w:1344

}else if msg,ok:=v.val.(*message);ok{


/*111:*/


//line amail.w:1418

{
glog.V(debug).Infof("appending a '%s' message to idmap\n",v.id)
if val,ok:=idmap[v.id];!ok{
glog.V(debug).Infof("'%s' message  doesn't exist, creating\n",v.id)
idmap[v.id]= &struct{msg*message;children idmessages}{msg,nil}
}else{
glog.V(debug).Infof("'%s' message exists, reseting\n",v.id)
val.msg= msg
for i,_:=range val.children{
val.children[i].parent= msg
}
idmap[v.id]= val
}
if len(msg.inreplyto)==0{
continue
}
if val,ok:=idmap[msg.inreplyto];!ok{
glog.V(debug).Infof("'%s' message  doesn't exist, creating\n",msg.inreplyto)
idmap[msg.inreplyto]= &struct{msg*message;children idmessages}{nil,append(idmessages{},msg)}
}else{
glog.V(debug).Infof("'%s' message exists, appending a child\n",msg.inreplyto)
if _,ok:=val.children.SearchInsert(msg);ok{
msg.parent= val.msg
}
}
}



/*:111*/


//line amail.w:1346

}else if ch,ok:=v.val.(chan idmessages);ok{


/*114:*/


//line amail.w:1491

{
if val,ok:=idmap[v.id];ok{
glog.V(debug).Infof("sending children for '%s'\n",v.id)
children:=make(idmessages,len(val.children),len(val.children))
copy(children,val.children)
sort.Sort(children)
ch<-children
}else{
glog.V(debug).Infof("'%s' is not found\n",v.id)
ch<-nil
}
}



/*:114*/


//line amail.w:1348

}



/*:102*/


//line amail.w:447

}
}




/*:33*/


//line amail.w:70

}



/*:2*/



/*24:*/


//line amail.w:256

func(this mailboxes)Len()int{
return len(this)
}

func(this mailboxes)Less(i,j int)bool{
return this[i].name<this[j].name
}

func(this mailboxes)Swap(i,j int){
t:=this[i]
this[i]= this[j]
this[j]= t
}



/*:24*/



/*30:*/


//line amail.w:319

func(this*mailbox)newMessage(id int)(msg*message,unread bool,err error){
glog.V(debug).Infof("newMessage: trying to open '%d/info'\n",id)
f,err:=this.fid.Walk(fmt.Sprintf("%d/info",id))
if err==nil{
err= f.Open(plan9.OREAD)
}
if err!=nil{
glog.Errorf("can't open to '%s/%d/info': %s\n",this.name,id,err)
return
}
defer f.Close()
msg= &message{id:id,box:this,

/*171:*/


//line amail.w:2309

cids:make(map[string]*file),



/*:171*/


//line amail.w:331
}
b:=bufio.NewReader(f)
unread= true
glog.V(debug).Infof("newMessage: reading and parsing of a content of '%d/info'\n",id)
for s,err:=b.ReadString('\n');err==nil;s,err= b.ReadString('\n'){
if s[len(s)-1]=='\n'{
s= s[:len(s)-1]
}
if strings.HasPrefix(s,"flags "){
if strings.Index(s,"seen")>=0{
unread= false
}
continue
}


/*92:*/


//line amail.w:1216

if strings.HasPrefix(s,"from "){
msg.from= s[len("from "):]
msg.from= strings.Replace(msg.from,"'' ","",-1)
continue
}
var unixdate int64
if _,err:=fmt.Sscanf(s,"unixdate %d",&unixdate);err==nil{
msg.date= time.Unix(unixdate,0)
continue
}
if strings.HasPrefix(s,"subject "){
msg.subject= s[len("subject "):]
continue
}




/*:92*/



/*101:*/


//line amail.w:1327

{
if _,err:=fmt.Sscanf(s,"inreplyto %s",&msg.inreplyto);err==nil{
msg.inreplyto= strings.Trim(msg.inreplyto,"<>")
continue
}
if _,err:=fmt.Sscanf(s,"messageid %s",&msg.messageid);err==nil{
msg.messageid= strings.Trim(msg.messageid,"<>")
idch<-struct{id string;val interface{}}{msg.messageid,msg}
continue
}
}



/*:101*/



/*160:*/


//line amail.w:2093

if strings.HasPrefix(s,"to "){
msg.to= split(s[len("to "):])
continue
}
if strings.HasPrefix(s,"cc "){
msg.cc= split(s[len("cc "):])
continue
}




/*:160*/


//line amail.w:345

}
msg.unread= unread
return

}



/*:30*/



/*36:*/


//line amail.w:483

func(this mailboxes)Search(name string)(int,bool){
pos:=sort.Search(len(this),
func(i int)bool{return this[i].name>=name});
if pos!=len(this)&&this[pos].name==name{
return pos,true
}
return pos,false
}



/*:36*/



/*52:*/


//line amail.w:736

func escape(s string)(res string){
for _,v:=range s{
if strings.ContainsRune("\\/[].+?()*^$",v){
res+= "\\"
}
res+= string(v)
}
return res
}



/*:52*/



/*55:*/


//line amail.w:762

func(this messages)Len()int{
return len(this)
}

func(this messages)Less(i,j int)bool{
return this[i].id> this[j].id
}

func(this messages)Swap(i,j int){
t:=this[i]
this[i]= this[j]
this[j]= t
}



/*:55*/



/*56:*/


//line amail.w:779

func(this messages)Search(id int)(int,bool){
pos:=sort.Search(len(this),func(i int)bool{return this[i].id<=id});
if pos!=len(this)&&this[pos].id==id{
return pos,true
}
return pos,false
}



/*:56*/



/*57:*/


//line amail.w:789

func(this*messages)Insert(msg*message,pos int){
*this= append(*this,nil)
copy((*this)[pos+1:],(*this)[pos:])
(*this)[pos]= msg
}



/*:57*/



/*58:*/


//line amail.w:799

func(this*messages)SearchInsert(msg*message)(int,bool){
pos,ok:=this.Search(msg.id)
if ok{
return pos,false
}
this.Insert(msg,pos)
return pos,true
}



/*:58*/



/*60:*/


//line amail.w:813

func(this*messages)Delete(pos int)(*message,bool){
if pos<0||pos> len(*this)-1{
return nil,false
}
msg:=(*this)[pos]
*this= append((*this)[:pos],(*this)[pos+1:]...)
return msg,true
}



/*:60*/



/*62:*/


//line amail.w:827

func(this*messages)DeleteById(id int)(*message,bool){
pos,ok:=this.Search(id)
if!ok{
return nil,false
}
return this.Delete(pos)
}



/*:62*/



/*64:*/


//line amail.w:841

func(box*mailbox)loop(){
glog.V(debug).Infof("start a message loop for the '%s' mailbox\n",box.name)
counted:=false
pcount:=0
ontop:=false


/*49:*/


//line amail.w:667

{
glog.V(debug).Infof("counting of messages in the '%s' mailbox\n",box.name)
var err error
box.fid,err= rfid.Walk(box.name)
if err!=nil{
glog.Errorf("can't walk to '%s': %v",box.name,err)
return
}
defer box.fid.Close()
fs,err:=box.fid.Dirreadall()
if err!=nil{
glog.Errorf("can't read a mailbox '%s': %s",box.name,err)
return
}
box.total= len(fs)-2
box.all= make(messages,0,box.total)
for i:=len(fs)-1;i>=0;{
select{


/*4:*/


//line amail.w:81

case<-exit:
glog.V(debug).Infoln("on exit!")
return



/*:4*/


//line amail.w:686



/*65:*/


//line amail.w:862

case id:=<-box.mch:
glog.V(debug).Infof("'%d' should be added to the '%s' mailbox\n",id,box.name)
msg,new,err:=box.newMessage(id)
if err!=nil{
continue
}
if new{


/*67:*/


//line amail.w:887

{
glog.V(debug).Infof("addition of the '%d' message to the list of unread messages of the '%s' mailbox\n",
msg.id,box.name)
box.unread.SearchInsert(msg)
}



/*:67*/


//line amail.w:870

}
box.total++


/*68:*/


//line amail.w:895

{
glog.V(debug).Infof("addition of the '%d' message to the list of all messages of the '%s' mailbox\n",
msg.id,box.name)
box.all.SearchInsert(msg)
}



/*:68*/


//line amail.w:873



/*126:*/


//line amail.w:1628

{
glog.V(debug).Infof("inform the '%s' mailbox to print a message '%d'\n",box.name,msg.id)
box.rfch<-&struct{seek bool;msgs messages}{true,append(messages{},msg)}
}



/*:126*/


//line amail.w:874



/*72:*/


//line amail.w:932

glog.V(debug).Infof("sending a snapshot of the '%s' mailbox to refresh the main window\n",box.name)
b:=*box
rfch<-&b



/*:72*/


//line amail.w:875

case id:=<-box.dch:
glog.V(debug).Infof("'%d' should be deleted from the '%s' mailbox\n",id,box.name)


/*69:*/


//line amail.w:903

box.deleteMessage(id)



/*:69*/


//line amail.w:878



/*155:*/


//line amail.w:2036

glog.V(debug).Infof("box.deleted:%d\n",box.deleted)
if box.deleted==0{
w:=box.w


/*142:*/


//line amail.w:1815

glog.V(debug).Infof("write a tag of the '%s' mailbox's window\n",box.name)
if err:=writeTag(box.w,fmt.Sprintf(" %sMail Delmesg %s%s %s Search ",
func()string{
if box.deleted> 0{
return"Put "
}
return""
}(),
func()string{
if box.thread{
if box.shownew{
return"ShowNew "
}else{
return"ShowAll "
}
}else if box.shownew||!box.showthreads{
return"Thread "
}
return""
}(),
func()string{if box.shownew{return"ShowAll"}else{return"ShowNew"}}(),
func()string{if box.showthreads{return"ShowPlain"}else{return"ShowThreads"}}()))
err!=nil{
glog.Errorf("can't set a tag of the '%s' box's window: %v\n",box.name,err)
}



/*:142*/


//line amail.w:2040



/*50:*/


//line amail.w:718

if w!=nil{
glog.V(debug).Infoln("setting the window to clean state")
if err:=w.WriteCtl("clean");err!=nil{
glog.Errorf("can't write to 'ctl' file: %s\n",err)
}
}



/*:50*/


//line amail.w:2041

}else{
w:=box.w


/*142:*/


//line amail.w:1815

glog.V(debug).Infof("write a tag of the '%s' mailbox's window\n",box.name)
if err:=writeTag(box.w,fmt.Sprintf(" %sMail Delmesg %s%s %s Search ",
func()string{
if box.deleted> 0{
return"Put "
}
return""
}(),
func()string{
if box.thread{
if box.shownew{
return"ShowNew "
}else{
return"ShowAll "
}
}else if box.shownew||!box.showthreads{
return"Thread "
}
return""
}(),
func()string{if box.shownew{return"ShowAll"}else{return"ShowNew"}}(),
func()string{if box.showthreads{return"ShowPlain"}else{return"ShowThreads"}}()))
err!=nil{
glog.Errorf("can't set a tag of the '%s' box's window: %v\n",box.name,err)
}



/*:142*/


//line amail.w:2044



/*51:*/


//line amail.w:727

if w!=nil{
glog.V(debug).Infoln("setting the window to dirty state")
if err:=w.WriteCtl("dirty");err!=nil{
glog.Errorf("can't write to 'ctl' file: %s\n",err)
}
}



/*:51*/


//line amail.w:2045

}



/*:155*/


//line amail.w:879



/*72:*/


//line amail.w:932

glog.V(debug).Infof("sending a snapshot of the '%s' mailbox to refresh the main window\n",box.name)
b:=*box
rfch<-&b



/*:72*/


//line amail.w:880




/*:65*/



/*77:*/


//line amail.w:969

case<-box.cch:
glog.V(debug).Infof("a signal to create the '%s' mailbox window has been received\n",box.name)
if box.w==nil{
box.shownew= shownew
box.showthreads= showthreads
box.thread= false


/*78:*/


//line amail.w:985

glog.V(debug).Infof("creation a window for the '%s' mailbox\n",box.name)
var err error
if box.w,err= goacme.New();err!=nil{
glog.Errorf("can't create a window: %v\n",err)
os.Exit(1)
}
if box.ech,err= box.w.EventChannel(0,goacme.Mouse,goacme.Look|goacme.Execute);err!=nil{
glog.Errorf("can't open an event channel of the window %v\n",err)
os.Exit(1)
}


/*141:*/


//line amail.w:1809

name:="Amail/"+box.name
w:=box.w


/*53:*/


//line amail.w:748

glog.V(debug).Infoln("printing a name for a window")
if err:=w.WriteCtl("name %s",strings.Replace(name," ","␣",-1));err!=nil{
glog.Errorf("can't write to 'ctl' file: %s\n",err)
}



/*:53*/


//line amail.w:1812




/*:141*/


//line amail.w:996



/*142:*/


//line amail.w:1815

glog.V(debug).Infof("write a tag of the '%s' mailbox's window\n",box.name)
if err:=writeTag(box.w,fmt.Sprintf(" %sMail Delmesg %s%s %s Search ",
func()string{
if box.deleted> 0{
return"Put "
}
return""
}(),
func()string{
if box.thread{
if box.shownew{
return"ShowNew "
}else{
return"ShowAll "
}
}else if box.shownew||!box.showthreads{
return"Thread "
}
return""
}(),
func()string{if box.shownew{return"ShowAll"}else{return"ShowNew"}}(),
func()string{if box.showthreads{return"ShowPlain"}else{return"ShowThreads"}}()))
err!=nil{
glog.Errorf("can't set a tag of the '%s' box's window: %v\n",box.name,err)
}



/*:142*/


//line amail.w:997





/*:78*/


//line amail.w:976



/*122:*/


//line amail.w:1574

{
glog.V(debug).Infof("inform the '%s' mailbox to print messages\n",box.name)


/*121:*/


//line amail.w:1564

var src messages
if box.shownew{
src= box.unread
}else{
src= box.all
}



/*:121*/


//line amail.w:1577

msgs:=append(messages{},src...)
if box.threadMode(){
src= msgs
msgs= make(messages,0,len(src))
for len(src)> 0{
msg:=src[0]


/*104:*/


//line amail.w:1356

for msg.parent!=nil{
msg= msg.parent
}



/*:104*/


//line amail.w:1584

glog.V(debug).Infof("root of thread: '%s/%d'\n",msg.box.name,msg.id)


/*123:*/


//line amail.w:1594

msgs= append(msgs,msg)
if p,ok:=src.Search(msg.id);ok&&src[p]==msg{
glog.V(debug).Infof("removing '%d' from src\n",src[p].id)
src.Delete(p)
}
msgs,src= getchildren(msg,msgs,src)




/*:123*/


//line amail.w:1586

}
}
box.rfch<-&struct{seek bool;msgs messages}{false,msgs}
}



/*:122*/


//line amail.w:977



/*8:*/


//line amail.w:109

glog.V(debug).Infoln("increase the windows count")
wch<-1



/*:8*/


//line amail.w:978

}else{
glog.V(debug).Infof("a window of the '%s' mailbox already exists, just show it\n",box.name)
box.w.WriteCtl("dot=addr\nshow")
}



/*:77*/



/*80:*/


//line amail.w:1005

case ev,ok:=<-box.ech:
glog.V(debug).Infof("an event has been received from the '%s' mailbox window: %v\n",box.name,ev)
if!ok{
box.ech= nil
continue
}
if(ev.Type&goacme.Execute)==goacme.Execute{
switch ev.Text{
case"Del":


/*129:*/


//line amail.w:1649

box.pos= 0
ontop= false



/*:129*/



/*135:*/


//line amail.w:1744

{
glog.V(debug).Infof("clean window-specific stuff of the '%s' mailbox\n",box.name)
close(box.irfch)
box.irfch= make(chan*struct{seek bool;msgs messages},100)
pcount= 0
ontop= false
}



/*:135*/


//line amail.w:1015

box.w.Del(true)
box.w.Close()
box.w= nil


/*9:*/


//line amail.w:114

glog.V(debug).Infoln("decrease the windows count")
wch<-(-1)




/*:9*/


//line amail.w:1019

continue
case"ShowNew":
box.thread= false
box.shownew= true
case"ShowAll":
box.thread= false
box.shownew= false
case"ShowThreads":
box.showthreads= true
if box.shownew==true{


/*142:*/


//line amail.w:1815

glog.V(debug).Infof("write a tag of the '%s' mailbox's window\n",box.name)
if err:=writeTag(box.w,fmt.Sprintf(" %sMail Delmesg %s%s %s Search ",
func()string{
if box.deleted> 0{
return"Put "
}
return""
}(),
func()string{
if box.thread{
if box.shownew{
return"ShowNew "
}else{
return"ShowAll "
}
}else if box.shownew||!box.showthreads{
return"Thread "
}
return""
}(),
func()string{if box.shownew{return"ShowAll"}else{return"ShowNew"}}(),
func()string{if box.showthreads{return"ShowPlain"}else{return"ShowThreads"}}()))
err!=nil{
glog.Errorf("can't set a tag of the '%s' box's window: %v\n",box.name,err)
}



/*:142*/


//line amail.w:1030

continue
}
case"ShowPlain":
box.showthreads= false
if box.shownew==true{


/*142:*/


//line amail.w:1815

glog.V(debug).Infof("write a tag of the '%s' mailbox's window\n",box.name)
if err:=writeTag(box.w,fmt.Sprintf(" %sMail Delmesg %s%s %s Search ",
func()string{
if box.deleted> 0{
return"Put "
}
return""
}(),
func()string{
if box.thread{
if box.shownew{
return"ShowNew "
}else{
return"ShowAll "
}
}else if box.shownew||!box.showthreads{
return"Thread "
}
return""
}(),
func()string{if box.shownew{return"ShowAll"}else{return"ShowNew"}}(),
func()string{if box.showthreads{return"ShowPlain"}else{return"ShowThreads"}}()))
err!=nil{
glog.Errorf("can't set a tag of the '%s' box's window: %v\n",box.name,err)
}



/*:142*/


//line amail.w:1036

continue
}
case"Thread":
var msg*message
if len(ev.Arg)==0{


/*88:*/


//line amail.w:1171

glog.V(debug).Infof("getting a pointer to current message in the window of the '%s' mailbox\n",box.name)
num:=0
if err:=box.w.WriteCtl("addr=dot");err!=nil{
glog.Errorf("can't write to 'ctl': %s\n",err)
}else if err:=box.w.WriteAddr("-/^/");err!=nil{
glog.V(debug).Infof("can't write to 'addr': %v\n",err)
}else if err:=box.w.WriteAddr("/[0-9]+(%s)?\\//",escape(newmark));err!=nil{
glog.V(debug).Infof("can't write to 'addr': %s\n",err)
}else if data,err:=box.w.File("data");err!=nil{
glog.Errorf("can't open 'data' file: %s\n",err)
}else if str,err:=bufio.NewReader(data).ReadString('/');err!=nil{
glog.Errorf("can't read from 'data' file: %s\n",err)
}else if _,err:=fmt.Sscanf(strings.TrimLeft(str,levelmark),"%d",&num);err==nil{
glog.V(debug).Infof("current message: %d(%s)\n",num,str)
if p,ok:=box.all.Search(num);ok{
msg= box.all[p]
}
}else{
glog.V(debug).Infof("can't get a current message from: %s\n",str)
}



/*:88*/


//line amail.w:1042

}else if num,err:=strconv.Atoi(strings.TrimSpace(ev.Arg));err!=nil{
continue
}else if p,ok:=box.all.Search(num);ok{
msg= box.all[p]
}
if msg!=nil{
box.thread= true


/*142:*/


//line amail.w:1815

glog.V(debug).Infof("write a tag of the '%s' mailbox's window\n",box.name)
if err:=writeTag(box.w,fmt.Sprintf(" %sMail Delmesg %s%s %s Search ",
func()string{
if box.deleted> 0{
return"Put "
}
return""
}(),
func()string{
if box.thread{
if box.shownew{
return"ShowNew "
}else{
return"ShowAll "
}
}else if box.shownew||!box.showthreads{
return"Thread "
}
return""
}(),
func()string{if box.shownew{return"ShowAll"}else{return"ShowNew"}}(),
func()string{if box.showthreads{return"ShowPlain"}else{return"ShowThreads"}}()))
err!=nil{
glog.Errorf("can't set a tag of the '%s' box's window: %v\n",box.name,err)
}



/*:142*/


//line amail.w:1050



/*143:*/


//line amail.w:1843

glog.V(debug).Infof("clean the '%s' mailbox's window\n",box.name)
clean(box.w)



/*:143*/


//line amail.w:1051



/*129:*/


//line amail.w:1649

box.pos= 0
ontop= false



/*:129*/



/*135:*/


//line amail.w:1744

{
glog.V(debug).Infof("clean window-specific stuff of the '%s' mailbox\n",box.name)
close(box.irfch)
box.irfch= make(chan*struct{seek bool;msgs messages},100)
pcount= 0
ontop= false
}



/*:135*/


//line amail.w:1052



/*125:*/


//line amail.w:1620



/*104:*/


//line amail.w:1356

for msg.parent!=nil{
msg= msg.parent
}



/*:104*/


//line amail.w:1621

var msgs messages
src:=append(messages{},msg)


/*123:*/


//line amail.w:1594

msgs= append(msgs,msg)
if p,ok:=src.Search(msg.id);ok&&src[p]==msg{
glog.V(debug).Infof("removing '%d' from src\n",src[p].id)
src.Delete(p)
}
msgs,src= getchildren(msg,msgs,src)




/*:123*/


//line amail.w:1624

box.rfch<-&struct{seek bool;msgs messages}{false,msgs}



/*:125*/


//line amail.w:1053

}
continue
case"Delmesg":


/*94:*/


//line amail.w:1240

if err:=box.w.WriteCtl("addr=dot");err!=nil{
glog.Errorf("can't write to 'ctl': %s\n",err)
}else if xdata,err:=box.w.File("xdata");err!=nil{
glog.Errorf("can't open 'xdata' file: %s\n",err)
}else{
b:=bufio.NewReader(xdata)
var msgs messages
for s,err:=b.ReadString('\n');err==nil||err==io.EOF;s,err= b.ReadString('\n'){
num:=0
glog.V(debug).Infof("looking a message number in '%s'\n",s)
if _,err:=fmt.Sscanf(strings.TrimLeft(s,levelmark+deleted),"%d",&num);err==nil{
glog.V(debug).Infof("the message number is '%d'\n",num)


/*95:*/


//line amail.w:1263

if p,ok:=box.all.Search(num);ok{
if box.all[p].deleted{
continue
}
box.all[p].deleted= true
box.deleted++
msgs= append(msgs,box.all[p])
if box.all[p].w!=nil{
this:=box.all[p]


/*163:*/


//line amail.w:2156

this.writeTag()



/*:163*/


//line amail.w:1273

}
}



/*:95*/


//line amail.w:1253

}
if err==io.EOF{
break
}
}


/*127:*/


//line amail.w:1635

{
if len(msgs)!=0{
glog.V(debug).Infof("inform the '%s' mailbox to print messages with setting a position\n",box.name)
box.rfch<-&struct{seek bool;msgs messages}{true,msgs}
}
}




/*:127*/


//line amail.w:1259

}



/*:94*/


//line amail.w:1057

continue
case"Put":


/*96:*/


//line amail.w:1278

f,err:=box.fid.Walk("ctl")
if err==nil{
err= f.Open(plan9.OWRITE)
}
if err!=nil{
glog.Errorf("can't open 'ctl': %v\n",err)
continue
}
for i:=0;i<len(box.all);{
if!box.all[i].deleted||box.all[i].w!=nil{
i++
continue
}
id:=box.all[i].id
glog.V(debug).Infof("deleting of '%s/%d' from the server\n",box.name,id)
if _,err:=f.Write([]byte(fmt.Sprintf("delete %s %d",box.name,id)));err!=nil{
glog.Errorf("can't delete the '%s/%d' message from the server : %v\n",box.name,id,err)
}


/*69:*/


//line amail.w:903

box.deleteMessage(id)



/*:69*/


//line amail.w:1297

}
f.Close()


/*72:*/


//line amail.w:932

glog.V(debug).Infof("sending a snapshot of the '%s' mailbox to refresh the main window\n",box.name)
b:=*box
rfch<-&b



/*:72*/


//line amail.w:1300



/*155:*/


//line amail.w:2036

glog.V(debug).Infof("box.deleted:%d\n",box.deleted)
if box.deleted==0{
w:=box.w


/*142:*/


//line amail.w:1815

glog.V(debug).Infof("write a tag of the '%s' mailbox's window\n",box.name)
if err:=writeTag(box.w,fmt.Sprintf(" %sMail Delmesg %s%s %s Search ",
func()string{
if box.deleted> 0{
return"Put "
}
return""
}(),
func()string{
if box.thread{
if box.shownew{
return"ShowNew "
}else{
return"ShowAll "
}
}else if box.shownew||!box.showthreads{
return"Thread "
}
return""
}(),
func()string{if box.shownew{return"ShowAll"}else{return"ShowNew"}}(),
func()string{if box.showthreads{return"ShowPlain"}else{return"ShowThreads"}}()))
err!=nil{
glog.Errorf("can't set a tag of the '%s' box's window: %v\n",box.name,err)
}



/*:142*/


//line amail.w:2040



/*50:*/


//line amail.w:718

if w!=nil{
glog.V(debug).Infoln("setting the window to clean state")
if err:=w.WriteCtl("clean");err!=nil{
glog.Errorf("can't write to 'ctl' file: %s\n",err)
}
}



/*:50*/


//line amail.w:2041

}else{
w:=box.w


/*142:*/


//line amail.w:1815

glog.V(debug).Infof("write a tag of the '%s' mailbox's window\n",box.name)
if err:=writeTag(box.w,fmt.Sprintf(" %sMail Delmesg %s%s %s Search ",
func()string{
if box.deleted> 0{
return"Put "
}
return""
}(),
func()string{
if box.thread{
if box.shownew{
return"ShowNew "
}else{
return"ShowAll "
}
}else if box.shownew||!box.showthreads{
return"Thread "
}
return""
}(),
func()string{if box.shownew{return"ShowAll"}else{return"ShowNew"}}(),
func()string{if box.showthreads{return"ShowPlain"}else{return"ShowThreads"}}()))
err!=nil{
glog.Errorf("can't set a tag of the '%s' box's window: %v\n",box.name,err)
}



/*:142*/


//line amail.w:2044



/*51:*/


//line amail.w:727

if w!=nil{
glog.V(debug).Infoln("setting the window to dirty state")
if err:=w.WriteCtl("dirty");err!=nil{
glog.Errorf("can't write to 'ctl' file: %s\n",err)
}
}



/*:51*/


//line amail.w:2045

}



/*:155*/


//line amail.w:1301




/*:96*/


//line amail.w:1060

continue
case"Mail":
var msg*message


/*192:*/


//line amail.w:2702

w,err:=goacme.New()
if err!=nil{
glog.Errorf("can't create a window: %v\n",err)
continue
}
if err:=writeTag(w," Look Post Undo ");err!=nil{
glog.Errorf("can't write a tag for a new message window: %v\n",err)
}


/*193:*/


//line amail.w:2714

go func(msg*message){
glog.V(debug).Infoln("starting a goroutine to process events from a composed mail window")
for ev,err:=w.ReadEvent();err==nil;ev,err= w.ReadEvent(){
if ev.Origin!=goacme.Mouse{
w.UnreadEvent(ev)
continue
}
if(ev.Type&goacme.Execute)==goacme.Execute{
switch ev.Text{
case"Del":
w.UnreadEvent(ev)
w.Close()
return
case"Post":


/*197:*/


//line amail.w:2801

{


/*199:*/


//line amail.w:2918



/*189:*/


//line amail.w:2633

once.Do(func(){

/*177:*/


//line amail.w:2450

if home= os.Getenv("home");len(home)==0{
if home= os.Getenv("HOME");len(home)==0{
glog.Errorln("can't get a home directory from the environment, the home is assumed '/'")
home= "/"
}
}




/*:177*/



/*186:*/


//line amail.w:2616

if u,err:=user.Current();err!=nil{
glog.Errorf("can't get a name of the current user: %v\n",err)
}else{
cuser= u.Username
}




/*:186*/



/*200:*/


//line amail.w:2922

if plan9dir= os.Getenv("PLAN9");len(plan9dir)==0{
glog.Errorln("can't get PLAN9 directory from the environment, the plan9dir is assumed '/usr/local/plan9'")
plan9dir= "/usr/local/plan9"
}



/*:200*/


//line amail.w:2634
})



/*:189*/


//line amail.w:2919




/*:199*/


//line amail.w:2803

w.Seek(0,0)
w.WriteAddr("0,$")
ff,_:=w.File("xdata")
b:=bufio.NewReader(ff)
var to,cc,bcc,attach,include[]string
var subject string
for{
s,err:=b.ReadString('\n')
if err!=nil{
break
}
s= strings.TrimSpace(s)
if len(s)==0{
// an empty line, the rest is a body of the message
break
}
p:=strings.Index(s,":")
if p!=-1{
f:=strings.Split(s[p+1:],",")
for i,_:=range f{
f[i]= strings.TrimSpace(f[i])
}
switch strings.ToLower(s[:p]){
case"from","to":
to= append(to,f...)
case"cc":
cc= append(cc,f...)
case"bcc":
bcc= append(bcc,f...)
case"attach":
attach= append(attach,f...)
case"include":
include= append(include,f...)
case"subject":
subject= fmt.Sprintf("%q",strings.TrimSpace(s[p+1:]))
}
}else{
// recipient addresses can be written without "to:"
f:=strings.Split(s,",")
for i,_:=range f{
f[i]= strings.TrimSpace(f[i])
}
to= append(to,f...)
}
}
args:=append([]string{},"-8")
if msg!=nil{
args= append(args,"-R",fmt.Sprintf("%s/%d",msg.box.name,msg.id))
}
if len(subject)!=0{
args= append(args,"-s",subject)
}
for _,v:=range include{
args= append(args,"-A",v)
}
for _,v:=range attach{
args= append(args,"-a",v)
}

c:=exec.Command(plan9dir+"/bin/upas/marshal",args...)
p,err:=c.StdinPipe()
if err!=nil{
glog.Errorf("can't get a stdin pipe: %v\n",err)
continue
}
if err:=c.Start();err!=nil{
glog.Errorf("can't start 'upas/marshal': %v\n",err)
continue
}
if len(to)!=0{
if _,err:=fmt.Fprintln(p,"To:",strings.Join(to,","));err!=nil{
glog.Errorf("can't write 'to' fields to 'upas/marshal': %v\n",err)
continue
}
}
glog.V(debug).Infoln("to is written")
if len(cc)!=0{
if _,err:=fmt.Fprintln(p,"CC:",strings.Join(cc,","));err!=nil{
glog.Errorf("can't write 'cc' fields to 'upas/marshal': %v\n",err)
continue
}
}
glog.V(debug).Infoln("cc is written")
if len(bcc)!=0{
if _,err:=fmt.Fprintln(p,"BCC:",strings.Join(bcc,","));err!=nil{
glog.Errorf("can't write 'bcc' fields to  'upas/marshal': %v\n",err)
continue
}
}
glog.V(debug).Infoln("bcc is written")
for s,err:=b.ReadString('\n');err==nil||err==io.EOF;s,err= b.ReadString('\n'){
glog.V(debug).Infof("writing '%s':%v",s,err)

p.Write([]byte(s))
if err==io.EOF{
p.Write([]byte("\n"))
break
}
}
glog.V(debug).Infoln("body is written")
p.Write([]byte("\n"))
p.Close()
c.Wait()
w.Del(true)
w.Close()
}





/*:197*/


//line amail.w:2729

}
}
w.UnreadEvent(ev)
}
}(msg)



/*:193*/


//line amail.w:2711




/*:192*/


//line amail.w:1064

name:=fmt.Sprintf("Amail/%s/New",box.name)


/*53:*/


//line amail.w:748

glog.V(debug).Infoln("printing a name for a window")
if err:=w.WriteCtl("name %s",strings.Replace(name," ","␣",-1));err!=nil{
glog.Errorf("can't write to 'ctl' file: %s\n",err)
}



/*:53*/


//line amail.w:1066

continue
case"Search":
glog.V(debug).Infof("search argument: '%s'\n",ev.Arg)


/*147:*/


//line amail.w:1931

{
msgs:=box.search(ev.Arg)


/*143:*/


//line amail.w:1843

glog.V(debug).Infof("clean the '%s' mailbox's window\n",box.name)
clean(box.w)



/*:143*/


//line amail.w:1934



/*129:*/


//line amail.w:1649

box.pos= 0
ontop= false



/*:129*/



/*135:*/


//line amail.w:1744

{
glog.V(debug).Infof("clean window-specific stuff of the '%s' mailbox\n",box.name)
close(box.irfch)
box.irfch= make(chan*struct{seek bool;msgs messages},100)
pcount= 0
ontop= false
}



/*:135*/


//line amail.w:1935

name:=fmt.Sprintf("Amail/%s/Search(%s)",box.name,strings.Replace(ev.Arg," ","␣",-1))
w:=box.w
box.thread= false
box.shownew= true


/*53:*/


//line amail.w:748

glog.V(debug).Infoln("printing a name for a window")
if err:=w.WriteCtl("name %s",strings.Replace(name," ","␣",-1));err!=nil{
glog.Errorf("can't write to 'ctl' file: %s\n",err)
}



/*:53*/


//line amail.w:1940

glog.V(debug).Infof("len of msgs: %v\n",len(msgs))
box.rfch<-&struct{seek bool;msgs messages}{false,msgs}
}



/*:147*/


//line amail.w:1070

continue
default:
box.w.UnreadEvent(ev)
continue
}


/*141:*/


//line amail.w:1809

name:="Amail/"+box.name
w:=box.w


/*53:*/


//line amail.w:748

glog.V(debug).Infoln("printing a name for a window")
if err:=w.WriteCtl("name %s",strings.Replace(name," ","␣",-1));err!=nil{
glog.Errorf("can't write to 'ctl' file: %s\n",err)
}



/*:53*/


//line amail.w:1812




/*:141*/


//line amail.w:1076



/*142:*/


//line amail.w:1815

glog.V(debug).Infof("write a tag of the '%s' mailbox's window\n",box.name)
if err:=writeTag(box.w,fmt.Sprintf(" %sMail Delmesg %s%s %s Search ",
func()string{
if box.deleted> 0{
return"Put "
}
return""
}(),
func()string{
if box.thread{
if box.shownew{
return"ShowNew "
}else{
return"ShowAll "
}
}else if box.shownew||!box.showthreads{
return"Thread "
}
return""
}(),
func()string{if box.shownew{return"ShowAll"}else{return"ShowNew"}}(),
func()string{if box.showthreads{return"ShowPlain"}else{return"ShowThreads"}}()))
err!=nil{
glog.Errorf("can't set a tag of the '%s' box's window: %v\n",box.name,err)
}



/*:142*/


//line amail.w:1077



/*143:*/


//line amail.w:1843

glog.V(debug).Infof("clean the '%s' mailbox's window\n",box.name)
clean(box.w)



/*:143*/


//line amail.w:1078



/*50:*/


//line amail.w:718

if w!=nil{
glog.V(debug).Infoln("setting the window to clean state")
if err:=w.WriteCtl("clean");err!=nil{
glog.Errorf("can't write to 'ctl' file: %s\n",err)
}
}



/*:50*/


//line amail.w:1079



/*129:*/


//line amail.w:1649

box.pos= 0
ontop= false



/*:129*/



/*135:*/


//line amail.w:1744

{
glog.V(debug).Infof("clean window-specific stuff of the '%s' mailbox\n",box.name)
close(box.irfch)
box.irfch= make(chan*struct{seek bool;msgs messages},100)
pcount= 0
ontop= false
}



/*:135*/


//line amail.w:1080



/*122:*/


//line amail.w:1574

{
glog.V(debug).Infof("inform the '%s' mailbox to print messages\n",box.name)


/*121:*/


//line amail.w:1564

var src messages
if box.shownew{
src= box.unread
}else{
src= box.all
}



/*:121*/


//line amail.w:1577

msgs:=append(messages{},src...)
if box.threadMode(){
src= msgs
msgs= make(messages,0,len(src))
for len(src)> 0{
msg:=src[0]


/*104:*/


//line amail.w:1356

for msg.parent!=nil{
msg= msg.parent
}



/*:104*/


//line amail.w:1584

glog.V(debug).Infof("root of thread: '%s/%d'\n",msg.box.name,msg.id)


/*123:*/


//line amail.w:1594

msgs= append(msgs,msg)
if p,ok:=src.Search(msg.id);ok&&src[p]==msg{
glog.V(debug).Infof("removing '%d' from src\n",src[p].id)
src.Delete(p)
}
msgs,src= getchildren(msg,msgs,src)




/*:123*/


//line amail.w:1586

}
}
box.rfch<-&struct{seek bool;msgs messages}{false,msgs}
}



/*:122*/


//line amail.w:1081

continue
}else if(ev.Type&goacme.Look)==goacme.Look{


/*84:*/


//line amail.w:1147

msgs:=make(map[string][]int)



/*:84*/


//line amail.w:1084

if(ev.Type&goacme.Tag)==goacme.Tag{
s:=ev.Text


/*82:*/


//line amail.w:1121

{
glog.V(debug).Infof("looking a message number in '%s'\n",s)
s= strings.TrimLeft(s,levelmark+deleted)
f:=strings.Split(s,"/")
glog.V(debug).Infof("parts of message path: '%v'\n",f)
num:=0
for i,v:=range f{
var err error
if num,err= strconv.Atoi(strings.TrimRight(v,newmark));err==nil{
name:=box.name
if i> 0{
name= strings.Join(f[:i],"/")
glog.V(debug).Infof("the message number is '%d' in the '%s' mailbox\n",num,name)
}


/*85:*/


//line amail.w:1151

glog.V(debug).Infof("sending a signal to open a window with the '%d' message of the '%s' mailbox\n",num,name)
msgs[name]= append(msgs[name],num)



/*:85*/


//line amail.w:1136

break
}
}
}



/*:82*/


//line amail.w:1087

}else{


/*81:*/


//line amail.w:1102

glog.V(debug).Infof("event: %v\n",ev)
if err:=box.w.WriteAddr("#%d,#%d",ev.Begin,ev.End);err!=nil{
glog.Errorf("can't write to 'addr': %s\n",err)
}else if xdata,err:=box.w.File("xdata");err!=nil{
glog.Errorf("can't open 'xdata' file: %s\n",err)
}else{
b:=bufio.NewReader(xdata)
for s,err:=b.ReadString('\n');err==nil||err==io.EOF;s,err= b.ReadString('\n'){


/*82:*/


//line amail.w:1121

{
glog.V(debug).Infof("looking a message number in '%s'\n",s)
s= strings.TrimLeft(s,levelmark+deleted)
f:=strings.Split(s,"/")
glog.V(debug).Infof("parts of message path: '%v'\n",f)
num:=0
for i,v:=range f{
var err error
if num,err= strconv.Atoi(strings.TrimRight(v,newmark));err==nil{
name:=box.name
if i> 0{
name= strings.Join(f[:i],"/")
glog.V(debug).Infof("the message number is '%d' in the '%s' mailbox\n",num,name)
}


/*85:*/


//line amail.w:1151

glog.V(debug).Infof("sending a signal to open a window with the '%d' message of the '%s' mailbox\n",num,name)
msgs[name]= append(msgs[name],num)



/*:85*/


//line amail.w:1136

break
}
}
}



/*:82*/


//line amail.w:1111

if err==io.EOF{
break
}
}

}



/*:81*/


//line amail.w:1089

}
if len(msgs)!=0{


/*87:*/


//line amail.w:1167

lch<-&msgs



/*:87*/


//line amail.w:1092

continue
}
}
box.w.UnreadEvent(ev)



/*:80*/



/*120:*/


//line amail.w:1546

case v:=<-box.rfch:
box.irfch<-v

case v:=<-box.irfch:
glog.V(debug).Infof("a signal to print message of the '%s' mailbox window has been received\n",box.name)
if box.w==nil{
glog.V(debug).Infof("a window of the '%s' mailbox doesn't exist, ignore the signal\n",box.name)
continue
}
if box.threadMode()&&!counted{
glog.V(debug).Infof("counting of threads of the '%s' mailbox is not finished, ignore the signal\n",box.name)
continue
}


/*132:*/


//line amail.w:1679

{
glog.V(debug).Infof("printing of messages of the '%s' mailbox from v.msgs, len(v.msgs): %d, with seeking a position: %v\n",box.name,len(v.msgs),v.seek)
f,err:=box.w.File("data")
if err!=nil{
glog.Errorf("can't open 'data' file of the '%s' messgebox: %v\n",box.name,err)
continue
}
if v.seek{


/*142:*/


//line amail.w:1815

glog.V(debug).Infof("write a tag of the '%s' mailbox's window\n",box.name)
if err:=writeTag(box.w,fmt.Sprintf(" %sMail Delmesg %s%s %s Search ",
func()string{
if box.deleted> 0{
return"Put "
}
return""
}(),
func()string{
if box.thread{
if box.shownew{
return"ShowNew "
}else{
return"ShowAll "
}
}else if box.shownew||!box.showthreads{
return"Thread "
}
return""
}(),
func()string{if box.shownew{return"ShowAll"}else{return"ShowNew"}}(),
func()string{if box.showthreads{return"ShowPlain"}else{return"ShowThreads"}}()))
err!=nil{
glog.Errorf("can't set a tag of the '%s' box's window: %v\n",box.name,err)
}



/*:142*/


//line amail.w:1688

msg:=v.msgs[0]


/*145:*/


//line amail.w:1861



/*121:*/


//line amail.w:1564

var src messages
if box.shownew{
src= box.unread
}else{
src= box.all
}



/*:121*/


//line amail.w:1862

addr:=fmt.Sprintf("0/^[%s]*(%s)?%s%d(%s)?\\/.*\\n\t.*\\n/",
escape(levelmark),
escape(deleted),
func()string{if box.name!=msg.box.name{return msg.box.name+"/"};return""}(),
msg.id,
escape(newmark))
glog.V(debug).Infoln("refreshed message addr:",addr)
if err:=box.w.WriteAddr(addr);err!=nil{
glog.V(debug).Infof("the '%d' message is not found in the window\n",msg.id)
if box.threadMode(){


/*146:*/


//line amail.w:1897

if msg.parent!=nil{
glog.V(debug).Infof("msg '%d' has a parent, looking for a position to print\n",msg.id)
m:=msg
msg:=m.parent
found:=false
for!found{


/*115:*/


//line amail.w:1506

ch:=make(chan idmessages)
glog.V(debug).Infof("getting children for '%s'\n",msg.messageid)
idch<-struct{id string;val interface{}}{msg.messageid,ch}
children:=<-ch



/*:115*/


//line amail.w:1904

if children==nil{
break
}
for i,v:=range children{
if v==m{
if i==0{
found= true
}
break
}
msg= v
}
}
if err:=box.w.WriteAddr("0/^[%s]*%s%s%d(%s)?\\/.*\\n\t.*\\n/+#0",
escape(levelmark),
func()string{if msg.deleted{return escape(deleted)};return""}(),
func()string{if box.name!=msg.box.name{return msg.box.name+"/"};return""}(),
msg.id,
escape(newmark));err!=nil{
glog.V(debug).Infof("can't write to 'addr': %s\n",err)
}
}else if err:=box.w.WriteAddr("#0-");err!=nil{
glog.Errorf("can't write to 'addr' file: %s\n",err)
}



/*:146*/


//line amail.w:1873

}else if p,ok:=src.Search(msg.id);!ok{
glog.V(debug).Infof("the '%d' message is not found\n",msg.id)
}else if p==0{
if err:=box.w.WriteAddr("#0-");err!=nil{
glog.Errorf("can't write to 'addr' file: %s\n",err)
}
}else if p==len(src)-1{
if err:=box.w.WriteAddr("$");err!=nil{
glog.Errorf("can't write to 'addr' file: %s\n",err)
}
}else if err:=box.w.WriteAddr("0/^[%s]*%s%s%d(%s)?\\/.*\\n\t.*\\n/+#0",
escape(levelmark),
func()string{if src[p-1].deleted{return escape(deleted)};return""}(),
func()string{if box.name!=src[p-1].box.name{return src[p-1].box.name+"/"};return""}(),
src[p-1].id,
escape(newmark));err!=nil{
glog.V(debug).Infof("can't write to 'addr': %s\n",err)
}
}



/*:145*/


//line amail.w:1690

}else if err:=box.w.WriteAddr("$");err!=nil{
glog.Errorf("can't write to 'addr' file: %s\n",err)
continue
}
w:=box.w
glog.V(debug).Infof("printing of messages of the '%s' mailbox\n",box.name)
buf:=make([]byte,0,0x8000)


/*134:*/


//line amail.w:1726

c:=0
for len(v.msgs)> 0&&c<100{
msg:=v.msgs[0]
glog.V(debug).Infof("printing of message with id: %v\n",msg.id)
if box.threadMode(){


/*136:*/


//line amail.w:1754

{
for p:=msg.parent;p!=nil;p= p.parent{
buf= append(buf,levelmark...)
}
}



/*:136*/


//line amail.w:1732

}
c++


/*90:*/


//line amail.w:1198

glog.V(debug).Infof("compose a header of the '%d' message of the '%s' mailbox\n",msg.id,box.name)
buf= append(buf,fmt.Sprintf("%s%s%d%s/\t%s\t%s\n\t%s\n",
func()string{if msg.deleted{return deleted};return""}(),
func()string{if msg.box.name!=box.name{return fmt.Sprintf("%s/",msg.box.name)};return""}(),
msg.id,
func()string{if msg.unread{return newmark};return""}(),
msg.from,
msg.date,
msg.subject)...)



/*:90*/


//line amail.w:1735

v.msgs= v.msgs[1:]
if v.seek{
break
}
}
pcount+= c



/*:134*/


//line amail.w:1698

if _,err:=f.Write(buf);err!=nil{
glog.Errorf("can't write to 'data' file of the '%s' messgebox: %v\n",box.name,err)
}


/*133:*/


//line amail.w:1712

if!ontop{
glog.V(debug).Infof("pcount: %v, ontop: %v\n",pcount,ontop)


/*93:*/


//line amail.w:1234

glog.V(debug).Infoln("Go to top of the window")
w.WriteAddr("#0")
w.WriteCtl("dot=addr\nshow")



/*:93*/


//line amail.w:1715

if pcount>=100{
ontop= true
}
}




/*:133*/


//line amail.w:1702

if len(v.msgs)> 0{
box.rfch<-&struct{seek bool;msgs messages}{v.seek,v.msgs}
}else{


/*155:*/


//line amail.w:2036

glog.V(debug).Infof("box.deleted:%d\n",box.deleted)
if box.deleted==0{
w:=box.w


/*142:*/


//line amail.w:1815

glog.V(debug).Infof("write a tag of the '%s' mailbox's window\n",box.name)
if err:=writeTag(box.w,fmt.Sprintf(" %sMail Delmesg %s%s %s Search ",
func()string{
if box.deleted> 0{
return"Put "
}
return""
}(),
func()string{
if box.thread{
if box.shownew{
return"ShowNew "
}else{
return"ShowAll "
}
}else if box.shownew||!box.showthreads{
return"Thread "
}
return""
}(),
func()string{if box.shownew{return"ShowAll"}else{return"ShowNew"}}(),
func()string{if box.showthreads{return"ShowPlain"}else{return"ShowThreads"}}()))
err!=nil{
glog.Errorf("can't set a tag of the '%s' box's window: %v\n",box.name,err)
}



/*:142*/


//line amail.w:2040



/*50:*/


//line amail.w:718

if w!=nil{
glog.V(debug).Infoln("setting the window to clean state")
if err:=w.WriteCtl("clean");err!=nil{
glog.Errorf("can't write to 'ctl' file: %s\n",err)
}
}



/*:50*/


//line amail.w:2041

}else{
w:=box.w


/*142:*/


//line amail.w:1815

glog.V(debug).Infof("write a tag of the '%s' mailbox's window\n",box.name)
if err:=writeTag(box.w,fmt.Sprintf(" %sMail Delmesg %s%s %s Search ",
func()string{
if box.deleted> 0{
return"Put "
}
return""
}(),
func()string{
if box.thread{
if box.shownew{
return"ShowNew "
}else{
return"ShowAll "
}
}else if box.shownew||!box.showthreads{
return"Thread "
}
return""
}(),
func()string{if box.shownew{return"ShowAll"}else{return"ShowNew"}}(),
func()string{if box.showthreads{return"ShowPlain"}else{return"ShowThreads"}}()))
err!=nil{
glog.Errorf("can't set a tag of the '%s' box's window: %v\n",box.name,err)
}



/*:142*/


//line amail.w:2044



/*51:*/


//line amail.w:727

if w!=nil{
glog.V(debug).Infoln("setting the window to dirty state")
if err:=w.WriteCtl("dirty");err!=nil{
glog.Errorf("can't write to 'ctl' file: %s\n",err)
}
}



/*:51*/


//line amail.w:2045

}



/*:155*/


//line amail.w:1706

}
}



/*:132*/


//line amail.w:1560





/*:120*/



/*152:*/


//line amail.w:1995

case ids:=<-box.lch:
for _,id:=range ids{
glog.V(debug).Infof("opening a window with the '%d' message of the '%s' mailbox\n",id,box.name)
p,ok:=box.all.Search(id)
if!ok{
glog.V(debug).Infof("the '%d' message of the '%s' mailbox has not found\n",id,box.name)
continue
}
msg:=box.all[p]
if msg.w==nil{
if err:=msg.open();err!=nil{
continue
}
if msg.unread{


/*153:*/


//line amail.w:2021

msg.unread= false
box.unread.DeleteById(id)



/*:153*/


//line amail.w:2010



/*154:*/


//line amail.w:2026

if box.thread||!box.shownew{


/*156:*/


//line amail.w:2049

if box.w!=nil{
glog.V(debug).Infof("removing the '%s' newmark from the '%d' message\n",newmark,id)
addr:=fmt.Sprintf("-/^[%s]*(%s)?%d/",escape(levelmark),escape(deleted),id)
if err:=box.w.WriteAddr(addr);err!=nil{
glog.Errorf("can't write '%s' to 'addr': %s\n",addr,err)
}else if err:=box.w.WriteAddr("/%s/",escape(newmark));err!=nil{
glog.Errorf("can't write to 'addr': %s\n",err)
}else if data,err:=box.w.File("data");err!=nil{
glog.Errorf("can't open 'data' file of the box '%s': %s\n",box.name,err)
}else if _,err:=data.Write([]byte{});err!=nil{
glog.Errorf("can't write to 'data' file of the box '%s': %s\n",box.name,err)
}
}



/*:156*/


//line amail.w:2028

}else{


/*157:*/


//line amail.w:2065

box.removeMessage(id)



/*:157*/


//line amail.w:2030

}


/*155:*/


//line amail.w:2036

glog.V(debug).Infof("box.deleted:%d\n",box.deleted)
if box.deleted==0{
w:=box.w


/*142:*/


//line amail.w:1815

glog.V(debug).Infof("write a tag of the '%s' mailbox's window\n",box.name)
if err:=writeTag(box.w,fmt.Sprintf(" %sMail Delmesg %s%s %s Search ",
func()string{
if box.deleted> 0{
return"Put "
}
return""
}(),
func()string{
if box.thread{
if box.shownew{
return"ShowNew "
}else{
return"ShowAll "
}
}else if box.shownew||!box.showthreads{
return"Thread "
}
return""
}(),
func()string{if box.shownew{return"ShowAll"}else{return"ShowNew"}}(),
func()string{if box.showthreads{return"ShowPlain"}else{return"ShowThreads"}}()))
err!=nil{
glog.Errorf("can't set a tag of the '%s' box's window: %v\n",box.name,err)
}



/*:142*/


//line amail.w:2040



/*50:*/


//line amail.w:718

if w!=nil{
glog.V(debug).Infoln("setting the window to clean state")
if err:=w.WriteCtl("clean");err!=nil{
glog.Errorf("can't write to 'ctl' file: %s\n",err)
}
}



/*:50*/


//line amail.w:2041

}else{
w:=box.w


/*142:*/


//line amail.w:1815

glog.V(debug).Infof("write a tag of the '%s' mailbox's window\n",box.name)
if err:=writeTag(box.w,fmt.Sprintf(" %sMail Delmesg %s%s %s Search ",
func()string{
if box.deleted> 0{
return"Put "
}
return""
}(),
func()string{
if box.thread{
if box.shownew{
return"ShowNew "
}else{
return"ShowAll "
}
}else if box.shownew||!box.showthreads{
return"Thread "
}
return""
}(),
func()string{if box.shownew{return"ShowAll"}else{return"ShowNew"}}(),
func()string{if box.showthreads{return"ShowPlain"}else{return"ShowThreads"}}()))
err!=nil{
glog.Errorf("can't set a tag of the '%s' box's window: %v\n",box.name,err)
}



/*:142*/


//line amail.w:2044



/*51:*/


//line amail.w:727

if w!=nil{
glog.V(debug).Infoln("setting the window to dirty state")
if err:=w.WriteCtl("dirty");err!=nil{
glog.Errorf("can't write to 'ctl' file: %s\n",err)
}
}



/*:51*/


//line amail.w:2045

}



/*:155*/


//line amail.w:2032





/*:154*/


//line amail.w:2011



/*72:*/


//line amail.w:932

glog.V(debug).Infof("sending a snapshot of the '%s' mailbox to refresh the main window\n",box.name)
b:=*box
rfch<-&b



/*:72*/


//line amail.w:2012

}
}else{
glog.V(debug).Infof("a window of the '%d' message of the '%s' already exists, just show it\n",id,box.name)
msg.w.WriteCtl("dot=addr\nshow")
}
}



/*:152*/


//line amail.w:687

default:
d:=fs[i]
i--
if(d.Mode&plan9.DMDIR)!=plan9.DMDIR{
continue
}
id,err:=strconv.Atoi(d.Name)
if err!=nil{// it seems this is a mailbox
// decrease a total number of messages
box.total--
name:=box.name+"/"+d.Name


/*27:*/


//line amail.w:303

glog.V(debug).Infof("send a mailbox '%s' to put in the list\n",name)
bch<-name



/*:27*/


//line amail.w:699

continue
}
if msg,new,_:=box.newMessage(id);err==nil{
if new{


/*67:*/


//line amail.w:887

{
glog.V(debug).Infof("addition of the '%d' message to the list of unread messages of the '%s' mailbox\n",
msg.id,box.name)
box.unread.SearchInsert(msg)
}



/*:67*/


//line amail.w:704

}


/*68:*/


//line amail.w:895

{
glog.V(debug).Infof("addition of the '%d' message to the list of all messages of the '%s' mailbox\n",
msg.id,box.name)
box.all.SearchInsert(msg)
}



/*:68*/


//line amail.w:706

}else{
glog.Errorf("can't create a new '%d' message in the '%s' mailbox\n",id,box.name)
}


/*72:*/


//line amail.w:932

glog.V(debug).Infof("sending a snapshot of the '%s' mailbox to refresh the main window\n",box.name)
b:=*box
rfch<-&b



/*:72*/


//line amail.w:710



/*130:*/


//line amail.w:1655

if!box.threadMode(){


/*121:*/


//line amail.w:1564

var src messages
if box.shownew{
src= box.unread
}else{
src= box.all
}



/*:121*/


//line amail.w:1657

if len(src)!=0&&len(src)%500==0{
glog.V(debug).Infof("inform the '%s' mailbox to print the last %d messages\n",box.name,len(src)-box.pos)
msgs:=append(messages{},src[box.pos:len(src)]...)
box.pos= len(src)
box.rfch<-&struct{seek bool;msgs messages}{false,msgs}
}
}



/*:130*/


//line amail.w:711

}
}


/*131:*/


//line amail.w:1667

if!box.threadMode(){


/*121:*/


//line amail.w:1564

var src messages
if box.shownew{
src= box.unread
}else{
src= box.all
}



/*:121*/


//line amail.w:1669

if box.pos!=len(src){
glog.V(debug).Infof("inform the '%s' mailbox to print the last %d messages\n",box.name,len(src)-box.pos)
msgs:=append(messages{},src[box.pos:len(src)]...)
box.pos= len(src)
box.rfch<-&struct{seek bool;msgs messages}{false,msgs}
}
}



/*:131*/


//line amail.w:714

}



/*:49*/


//line amail.w:847

counted= true
if box.threadMode(){


/*122:*/


//line amail.w:1574

{
glog.V(debug).Infof("inform the '%s' mailbox to print messages\n",box.name)


/*121:*/


//line amail.w:1564

var src messages
if box.shownew{
src= box.unread
}else{
src= box.all
}



/*:121*/


//line amail.w:1577

msgs:=append(messages{},src...)
if box.threadMode(){
src= msgs
msgs= make(messages,0,len(src))
for len(src)> 0{
msg:=src[0]


/*104:*/


//line amail.w:1356

for msg.parent!=nil{
msg= msg.parent
}



/*:104*/


//line amail.w:1584

glog.V(debug).Infof("root of thread: '%s/%d'\n",msg.box.name,msg.id)


/*123:*/


//line amail.w:1594

msgs= append(msgs,msg)
if p,ok:=src.Search(msg.id);ok&&src[p]==msg{
glog.V(debug).Infof("removing '%d' from src\n",src[p].id)
src.Delete(p)
}
msgs,src= getchildren(msg,msgs,src)




/*:123*/


//line amail.w:1586

}
}
box.rfch<-&struct{seek bool;msgs messages}{false,msgs}
}



/*:122*/


//line amail.w:850

}
defer glog.V(debug).Infof("a message loop of the '%s' mailbox is done\n",box.name)
for{
select{


/*4:*/


//line amail.w:81

case<-exit:
glog.V(debug).Infoln("on exit!")
return



/*:4*/


//line amail.w:855



/*65:*/


//line amail.w:862

case id:=<-box.mch:
glog.V(debug).Infof("'%d' should be added to the '%s' mailbox\n",id,box.name)
msg,new,err:=box.newMessage(id)
if err!=nil{
continue
}
if new{


/*67:*/


//line amail.w:887

{
glog.V(debug).Infof("addition of the '%d' message to the list of unread messages of the '%s' mailbox\n",
msg.id,box.name)
box.unread.SearchInsert(msg)
}



/*:67*/


//line amail.w:870

}
box.total++


/*68:*/


//line amail.w:895

{
glog.V(debug).Infof("addition of the '%d' message to the list of all messages of the '%s' mailbox\n",
msg.id,box.name)
box.all.SearchInsert(msg)
}



/*:68*/


//line amail.w:873



/*126:*/


//line amail.w:1628

{
glog.V(debug).Infof("inform the '%s' mailbox to print a message '%d'\n",box.name,msg.id)
box.rfch<-&struct{seek bool;msgs messages}{true,append(messages{},msg)}
}



/*:126*/


//line amail.w:874



/*72:*/


//line amail.w:932

glog.V(debug).Infof("sending a snapshot of the '%s' mailbox to refresh the main window\n",box.name)
b:=*box
rfch<-&b



/*:72*/


//line amail.w:875

case id:=<-box.dch:
glog.V(debug).Infof("'%d' should be deleted from the '%s' mailbox\n",id,box.name)


/*69:*/


//line amail.w:903

box.deleteMessage(id)



/*:69*/


//line amail.w:878



/*155:*/


//line amail.w:2036

glog.V(debug).Infof("box.deleted:%d\n",box.deleted)
if box.deleted==0{
w:=box.w


/*142:*/


//line amail.w:1815

glog.V(debug).Infof("write a tag of the '%s' mailbox's window\n",box.name)
if err:=writeTag(box.w,fmt.Sprintf(" %sMail Delmesg %s%s %s Search ",
func()string{
if box.deleted> 0{
return"Put "
}
return""
}(),
func()string{
if box.thread{
if box.shownew{
return"ShowNew "
}else{
return"ShowAll "
}
}else if box.shownew||!box.showthreads{
return"Thread "
}
return""
}(),
func()string{if box.shownew{return"ShowAll"}else{return"ShowNew"}}(),
func()string{if box.showthreads{return"ShowPlain"}else{return"ShowThreads"}}()))
err!=nil{
glog.Errorf("can't set a tag of the '%s' box's window: %v\n",box.name,err)
}



/*:142*/


//line amail.w:2040



/*50:*/


//line amail.w:718

if w!=nil{
glog.V(debug).Infoln("setting the window to clean state")
if err:=w.WriteCtl("clean");err!=nil{
glog.Errorf("can't write to 'ctl' file: %s\n",err)
}
}



/*:50*/


//line amail.w:2041

}else{
w:=box.w


/*142:*/


//line amail.w:1815

glog.V(debug).Infof("write a tag of the '%s' mailbox's window\n",box.name)
if err:=writeTag(box.w,fmt.Sprintf(" %sMail Delmesg %s%s %s Search ",
func()string{
if box.deleted> 0{
return"Put "
}
return""
}(),
func()string{
if box.thread{
if box.shownew{
return"ShowNew "
}else{
return"ShowAll "
}
}else if box.shownew||!box.showthreads{
return"Thread "
}
return""
}(),
func()string{if box.shownew{return"ShowAll"}else{return"ShowNew"}}(),
func()string{if box.showthreads{return"ShowPlain"}else{return"ShowThreads"}}()))
err!=nil{
glog.Errorf("can't set a tag of the '%s' box's window: %v\n",box.name,err)
}



/*:142*/


//line amail.w:2044



/*51:*/


//line amail.w:727

if w!=nil{
glog.V(debug).Infoln("setting the window to dirty state")
if err:=w.WriteCtl("dirty");err!=nil{
glog.Errorf("can't write to 'ctl' file: %s\n",err)
}
}



/*:51*/


//line amail.w:2045

}



/*:155*/


//line amail.w:879



/*72:*/


//line amail.w:932

glog.V(debug).Infof("sending a snapshot of the '%s' mailbox to refresh the main window\n",box.name)
b:=*box
rfch<-&b



/*:72*/


//line amail.w:880




/*:65*/



/*77:*/


//line amail.w:969

case<-box.cch:
glog.V(debug).Infof("a signal to create the '%s' mailbox window has been received\n",box.name)
if box.w==nil{
box.shownew= shownew
box.showthreads= showthreads
box.thread= false


/*78:*/


//line amail.w:985

glog.V(debug).Infof("creation a window for the '%s' mailbox\n",box.name)
var err error
if box.w,err= goacme.New();err!=nil{
glog.Errorf("can't create a window: %v\n",err)
os.Exit(1)
}
if box.ech,err= box.w.EventChannel(0,goacme.Mouse,goacme.Look|goacme.Execute);err!=nil{
glog.Errorf("can't open an event channel of the window %v\n",err)
os.Exit(1)
}


/*141:*/


//line amail.w:1809

name:="Amail/"+box.name
w:=box.w


/*53:*/


//line amail.w:748

glog.V(debug).Infoln("printing a name for a window")
if err:=w.WriteCtl("name %s",strings.Replace(name," ","␣",-1));err!=nil{
glog.Errorf("can't write to 'ctl' file: %s\n",err)
}



/*:53*/


//line amail.w:1812




/*:141*/


//line amail.w:996



/*142:*/


//line amail.w:1815

glog.V(debug).Infof("write a tag of the '%s' mailbox's window\n",box.name)
if err:=writeTag(box.w,fmt.Sprintf(" %sMail Delmesg %s%s %s Search ",
func()string{
if box.deleted> 0{
return"Put "
}
return""
}(),
func()string{
if box.thread{
if box.shownew{
return"ShowNew "
}else{
return"ShowAll "
}
}else if box.shownew||!box.showthreads{
return"Thread "
}
return""
}(),
func()string{if box.shownew{return"ShowAll"}else{return"ShowNew"}}(),
func()string{if box.showthreads{return"ShowPlain"}else{return"ShowThreads"}}()))
err!=nil{
glog.Errorf("can't set a tag of the '%s' box's window: %v\n",box.name,err)
}



/*:142*/


//line amail.w:997





/*:78*/


//line amail.w:976



/*122:*/


//line amail.w:1574

{
glog.V(debug).Infof("inform the '%s' mailbox to print messages\n",box.name)


/*121:*/


//line amail.w:1564

var src messages
if box.shownew{
src= box.unread
}else{
src= box.all
}



/*:121*/


//line amail.w:1577

msgs:=append(messages{},src...)
if box.threadMode(){
src= msgs
msgs= make(messages,0,len(src))
for len(src)> 0{
msg:=src[0]


/*104:*/


//line amail.w:1356

for msg.parent!=nil{
msg= msg.parent
}



/*:104*/


//line amail.w:1584

glog.V(debug).Infof("root of thread: '%s/%d'\n",msg.box.name,msg.id)


/*123:*/


//line amail.w:1594

msgs= append(msgs,msg)
if p,ok:=src.Search(msg.id);ok&&src[p]==msg{
glog.V(debug).Infof("removing '%d' from src\n",src[p].id)
src.Delete(p)
}
msgs,src= getchildren(msg,msgs,src)




/*:123*/


//line amail.w:1586

}
}
box.rfch<-&struct{seek bool;msgs messages}{false,msgs}
}



/*:122*/


//line amail.w:977



/*8:*/


//line amail.w:109

glog.V(debug).Infoln("increase the windows count")
wch<-1



/*:8*/


//line amail.w:978

}else{
glog.V(debug).Infof("a window of the '%s' mailbox already exists, just show it\n",box.name)
box.w.WriteCtl("dot=addr\nshow")
}



/*:77*/



/*80:*/


//line amail.w:1005

case ev,ok:=<-box.ech:
glog.V(debug).Infof("an event has been received from the '%s' mailbox window: %v\n",box.name,ev)
if!ok{
box.ech= nil
continue
}
if(ev.Type&goacme.Execute)==goacme.Execute{
switch ev.Text{
case"Del":


/*129:*/


//line amail.w:1649

box.pos= 0
ontop= false



/*:129*/



/*135:*/


//line amail.w:1744

{
glog.V(debug).Infof("clean window-specific stuff of the '%s' mailbox\n",box.name)
close(box.irfch)
box.irfch= make(chan*struct{seek bool;msgs messages},100)
pcount= 0
ontop= false
}



/*:135*/


//line amail.w:1015

box.w.Del(true)
box.w.Close()
box.w= nil


/*9:*/


//line amail.w:114

glog.V(debug).Infoln("decrease the windows count")
wch<-(-1)




/*:9*/


//line amail.w:1019

continue
case"ShowNew":
box.thread= false
box.shownew= true
case"ShowAll":
box.thread= false
box.shownew= false
case"ShowThreads":
box.showthreads= true
if box.shownew==true{


/*142:*/


//line amail.w:1815

glog.V(debug).Infof("write a tag of the '%s' mailbox's window\n",box.name)
if err:=writeTag(box.w,fmt.Sprintf(" %sMail Delmesg %s%s %s Search ",
func()string{
if box.deleted> 0{
return"Put "
}
return""
}(),
func()string{
if box.thread{
if box.shownew{
return"ShowNew "
}else{
return"ShowAll "
}
}else if box.shownew||!box.showthreads{
return"Thread "
}
return""
}(),
func()string{if box.shownew{return"ShowAll"}else{return"ShowNew"}}(),
func()string{if box.showthreads{return"ShowPlain"}else{return"ShowThreads"}}()))
err!=nil{
glog.Errorf("can't set a tag of the '%s' box's window: %v\n",box.name,err)
}



/*:142*/


//line amail.w:1030

continue
}
case"ShowPlain":
box.showthreads= false
if box.shownew==true{


/*142:*/


//line amail.w:1815

glog.V(debug).Infof("write a tag of the '%s' mailbox's window\n",box.name)
if err:=writeTag(box.w,fmt.Sprintf(" %sMail Delmesg %s%s %s Search ",
func()string{
if box.deleted> 0{
return"Put "
}
return""
}(),
func()string{
if box.thread{
if box.shownew{
return"ShowNew "
}else{
return"ShowAll "
}
}else if box.shownew||!box.showthreads{
return"Thread "
}
return""
}(),
func()string{if box.shownew{return"ShowAll"}else{return"ShowNew"}}(),
func()string{if box.showthreads{return"ShowPlain"}else{return"ShowThreads"}}()))
err!=nil{
glog.Errorf("can't set a tag of the '%s' box's window: %v\n",box.name,err)
}



/*:142*/


//line amail.w:1036

continue
}
case"Thread":
var msg*message
if len(ev.Arg)==0{


/*88:*/


//line amail.w:1171

glog.V(debug).Infof("getting a pointer to current message in the window of the '%s' mailbox\n",box.name)
num:=0
if err:=box.w.WriteCtl("addr=dot");err!=nil{
glog.Errorf("can't write to 'ctl': %s\n",err)
}else if err:=box.w.WriteAddr("-/^/");err!=nil{
glog.V(debug).Infof("can't write to 'addr': %v\n",err)
}else if err:=box.w.WriteAddr("/[0-9]+(%s)?\\//",escape(newmark));err!=nil{
glog.V(debug).Infof("can't write to 'addr': %s\n",err)
}else if data,err:=box.w.File("data");err!=nil{
glog.Errorf("can't open 'data' file: %s\n",err)
}else if str,err:=bufio.NewReader(data).ReadString('/');err!=nil{
glog.Errorf("can't read from 'data' file: %s\n",err)
}else if _,err:=fmt.Sscanf(strings.TrimLeft(str,levelmark),"%d",&num);err==nil{
glog.V(debug).Infof("current message: %d(%s)\n",num,str)
if p,ok:=box.all.Search(num);ok{
msg= box.all[p]
}
}else{
glog.V(debug).Infof("can't get a current message from: %s\n",str)
}



/*:88*/


//line amail.w:1042

}else if num,err:=strconv.Atoi(strings.TrimSpace(ev.Arg));err!=nil{
continue
}else if p,ok:=box.all.Search(num);ok{
msg= box.all[p]
}
if msg!=nil{
box.thread= true


/*142:*/


//line amail.w:1815

glog.V(debug).Infof("write a tag of the '%s' mailbox's window\n",box.name)
if err:=writeTag(box.w,fmt.Sprintf(" %sMail Delmesg %s%s %s Search ",
func()string{
if box.deleted> 0{
return"Put "
}
return""
}(),
func()string{
if box.thread{
if box.shownew{
return"ShowNew "
}else{
return"ShowAll "
}
}else if box.shownew||!box.showthreads{
return"Thread "
}
return""
}(),
func()string{if box.shownew{return"ShowAll"}else{return"ShowNew"}}(),
func()string{if box.showthreads{return"ShowPlain"}else{return"ShowThreads"}}()))
err!=nil{
glog.Errorf("can't set a tag of the '%s' box's window: %v\n",box.name,err)
}



/*:142*/


//line amail.w:1050



/*143:*/


//line amail.w:1843

glog.V(debug).Infof("clean the '%s' mailbox's window\n",box.name)
clean(box.w)



/*:143*/


//line amail.w:1051



/*129:*/


//line amail.w:1649

box.pos= 0
ontop= false



/*:129*/



/*135:*/


//line amail.w:1744

{
glog.V(debug).Infof("clean window-specific stuff of the '%s' mailbox\n",box.name)
close(box.irfch)
box.irfch= make(chan*struct{seek bool;msgs messages},100)
pcount= 0
ontop= false
}



/*:135*/


//line amail.w:1052



/*125:*/


//line amail.w:1620



/*104:*/


//line amail.w:1356

for msg.parent!=nil{
msg= msg.parent
}



/*:104*/


//line amail.w:1621

var msgs messages
src:=append(messages{},msg)


/*123:*/


//line amail.w:1594

msgs= append(msgs,msg)
if p,ok:=src.Search(msg.id);ok&&src[p]==msg{
glog.V(debug).Infof("removing '%d' from src\n",src[p].id)
src.Delete(p)
}
msgs,src= getchildren(msg,msgs,src)




/*:123*/


//line amail.w:1624

box.rfch<-&struct{seek bool;msgs messages}{false,msgs}



/*:125*/


//line amail.w:1053

}
continue
case"Delmesg":


/*94:*/


//line amail.w:1240

if err:=box.w.WriteCtl("addr=dot");err!=nil{
glog.Errorf("can't write to 'ctl': %s\n",err)
}else if xdata,err:=box.w.File("xdata");err!=nil{
glog.Errorf("can't open 'xdata' file: %s\n",err)
}else{
b:=bufio.NewReader(xdata)
var msgs messages
for s,err:=b.ReadString('\n');err==nil||err==io.EOF;s,err= b.ReadString('\n'){
num:=0
glog.V(debug).Infof("looking a message number in '%s'\n",s)
if _,err:=fmt.Sscanf(strings.TrimLeft(s,levelmark+deleted),"%d",&num);err==nil{
glog.V(debug).Infof("the message number is '%d'\n",num)


/*95:*/


//line amail.w:1263

if p,ok:=box.all.Search(num);ok{
if box.all[p].deleted{
continue
}
box.all[p].deleted= true
box.deleted++
msgs= append(msgs,box.all[p])
if box.all[p].w!=nil{
this:=box.all[p]


/*163:*/


//line amail.w:2156

this.writeTag()



/*:163*/


//line amail.w:1273

}
}



/*:95*/


//line amail.w:1253

}
if err==io.EOF{
break
}
}


/*127:*/


//line amail.w:1635

{
if len(msgs)!=0{
glog.V(debug).Infof("inform the '%s' mailbox to print messages with setting a position\n",box.name)
box.rfch<-&struct{seek bool;msgs messages}{true,msgs}
}
}




/*:127*/


//line amail.w:1259

}



/*:94*/


//line amail.w:1057

continue
case"Put":


/*96:*/


//line amail.w:1278

f,err:=box.fid.Walk("ctl")
if err==nil{
err= f.Open(plan9.OWRITE)
}
if err!=nil{
glog.Errorf("can't open 'ctl': %v\n",err)
continue
}
for i:=0;i<len(box.all);{
if!box.all[i].deleted||box.all[i].w!=nil{
i++
continue
}
id:=box.all[i].id
glog.V(debug).Infof("deleting of '%s/%d' from the server\n",box.name,id)
if _,err:=f.Write([]byte(fmt.Sprintf("delete %s %d",box.name,id)));err!=nil{
glog.Errorf("can't delete the '%s/%d' message from the server : %v\n",box.name,id,err)
}


/*69:*/


//line amail.w:903

box.deleteMessage(id)



/*:69*/


//line amail.w:1297

}
f.Close()


/*72:*/


//line amail.w:932

glog.V(debug).Infof("sending a snapshot of the '%s' mailbox to refresh the main window\n",box.name)
b:=*box
rfch<-&b



/*:72*/


//line amail.w:1300



/*155:*/


//line amail.w:2036

glog.V(debug).Infof("box.deleted:%d\n",box.deleted)
if box.deleted==0{
w:=box.w


/*142:*/


//line amail.w:1815

glog.V(debug).Infof("write a tag of the '%s' mailbox's window\n",box.name)
if err:=writeTag(box.w,fmt.Sprintf(" %sMail Delmesg %s%s %s Search ",
func()string{
if box.deleted> 0{
return"Put "
}
return""
}(),
func()string{
if box.thread{
if box.shownew{
return"ShowNew "
}else{
return"ShowAll "
}
}else if box.shownew||!box.showthreads{
return"Thread "
}
return""
}(),
func()string{if box.shownew{return"ShowAll"}else{return"ShowNew"}}(),
func()string{if box.showthreads{return"ShowPlain"}else{return"ShowThreads"}}()))
err!=nil{
glog.Errorf("can't set a tag of the '%s' box's window: %v\n",box.name,err)
}



/*:142*/


//line amail.w:2040



/*50:*/


//line amail.w:718

if w!=nil{
glog.V(debug).Infoln("setting the window to clean state")
if err:=w.WriteCtl("clean");err!=nil{
glog.Errorf("can't write to 'ctl' file: %s\n",err)
}
}



/*:50*/


//line amail.w:2041

}else{
w:=box.w


/*142:*/


//line amail.w:1815

glog.V(debug).Infof("write a tag of the '%s' mailbox's window\n",box.name)
if err:=writeTag(box.w,fmt.Sprintf(" %sMail Delmesg %s%s %s Search ",
func()string{
if box.deleted> 0{
return"Put "
}
return""
}(),
func()string{
if box.thread{
if box.shownew{
return"ShowNew "
}else{
return"ShowAll "
}
}else if box.shownew||!box.showthreads{
return"Thread "
}
return""
}(),
func()string{if box.shownew{return"ShowAll"}else{return"ShowNew"}}(),
func()string{if box.showthreads{return"ShowPlain"}else{return"ShowThreads"}}()))
err!=nil{
glog.Errorf("can't set a tag of the '%s' box's window: %v\n",box.name,err)
}



/*:142*/


//line amail.w:2044



/*51:*/


//line amail.w:727

if w!=nil{
glog.V(debug).Infoln("setting the window to dirty state")
if err:=w.WriteCtl("dirty");err!=nil{
glog.Errorf("can't write to 'ctl' file: %s\n",err)
}
}



/*:51*/


//line amail.w:2045

}



/*:155*/


//line amail.w:1301




/*:96*/


//line amail.w:1060

continue
case"Mail":
var msg*message


/*192:*/


//line amail.w:2702

w,err:=goacme.New()
if err!=nil{
glog.Errorf("can't create a window: %v\n",err)
continue
}
if err:=writeTag(w," Look Post Undo ");err!=nil{
glog.Errorf("can't write a tag for a new message window: %v\n",err)
}


/*193:*/


//line amail.w:2714

go func(msg*message){
glog.V(debug).Infoln("starting a goroutine to process events from a composed mail window")
for ev,err:=w.ReadEvent();err==nil;ev,err= w.ReadEvent(){
if ev.Origin!=goacme.Mouse{
w.UnreadEvent(ev)
continue
}
if(ev.Type&goacme.Execute)==goacme.Execute{
switch ev.Text{
case"Del":
w.UnreadEvent(ev)
w.Close()
return
case"Post":


/*197:*/


//line amail.w:2801

{


/*199:*/


//line amail.w:2918



/*189:*/


//line amail.w:2633

once.Do(func(){

/*177:*/


//line amail.w:2450

if home= os.Getenv("home");len(home)==0{
if home= os.Getenv("HOME");len(home)==0{
glog.Errorln("can't get a home directory from the environment, the home is assumed '/'")
home= "/"
}
}




/*:177*/



/*186:*/


//line amail.w:2616

if u,err:=user.Current();err!=nil{
glog.Errorf("can't get a name of the current user: %v\n",err)
}else{
cuser= u.Username
}




/*:186*/



/*200:*/


//line amail.w:2922

if plan9dir= os.Getenv("PLAN9");len(plan9dir)==0{
glog.Errorln("can't get PLAN9 directory from the environment, the plan9dir is assumed '/usr/local/plan9'")
plan9dir= "/usr/local/plan9"
}



/*:200*/


//line amail.w:2634
})



/*:189*/


//line amail.w:2919




/*:199*/


//line amail.w:2803

w.Seek(0,0)
w.WriteAddr("0,$")
ff,_:=w.File("xdata")
b:=bufio.NewReader(ff)
var to,cc,bcc,attach,include[]string
var subject string
for{
s,err:=b.ReadString('\n')
if err!=nil{
break
}
s= strings.TrimSpace(s)
if len(s)==0{
// an empty line, the rest is a body of the message
break
}
p:=strings.Index(s,":")
if p!=-1{
f:=strings.Split(s[p+1:],",")
for i,_:=range f{
f[i]= strings.TrimSpace(f[i])
}
switch strings.ToLower(s[:p]){
case"from","to":
to= append(to,f...)
case"cc":
cc= append(cc,f...)
case"bcc":
bcc= append(bcc,f...)
case"attach":
attach= append(attach,f...)
case"include":
include= append(include,f...)
case"subject":
subject= fmt.Sprintf("%q",strings.TrimSpace(s[p+1:]))
}
}else{
// recipient addresses can be written without "to:"
f:=strings.Split(s,",")
for i,_:=range f{
f[i]= strings.TrimSpace(f[i])
}
to= append(to,f...)
}
}
args:=append([]string{},"-8")
if msg!=nil{
args= append(args,"-R",fmt.Sprintf("%s/%d",msg.box.name,msg.id))
}
if len(subject)!=0{
args= append(args,"-s",subject)
}
for _,v:=range include{
args= append(args,"-A",v)
}
for _,v:=range attach{
args= append(args,"-a",v)
}

c:=exec.Command(plan9dir+"/bin/upas/marshal",args...)
p,err:=c.StdinPipe()
if err!=nil{
glog.Errorf("can't get a stdin pipe: %v\n",err)
continue
}
if err:=c.Start();err!=nil{
glog.Errorf("can't start 'upas/marshal': %v\n",err)
continue
}
if len(to)!=0{
if _,err:=fmt.Fprintln(p,"To:",strings.Join(to,","));err!=nil{
glog.Errorf("can't write 'to' fields to 'upas/marshal': %v\n",err)
continue
}
}
glog.V(debug).Infoln("to is written")
if len(cc)!=0{
if _,err:=fmt.Fprintln(p,"CC:",strings.Join(cc,","));err!=nil{
glog.Errorf("can't write 'cc' fields to 'upas/marshal': %v\n",err)
continue
}
}
glog.V(debug).Infoln("cc is written")
if len(bcc)!=0{
if _,err:=fmt.Fprintln(p,"BCC:",strings.Join(bcc,","));err!=nil{
glog.Errorf("can't write 'bcc' fields to  'upas/marshal': %v\n",err)
continue
}
}
glog.V(debug).Infoln("bcc is written")
for s,err:=b.ReadString('\n');err==nil||err==io.EOF;s,err= b.ReadString('\n'){
glog.V(debug).Infof("writing '%s':%v",s,err)

p.Write([]byte(s))
if err==io.EOF{
p.Write([]byte("\n"))
break
}
}
glog.V(debug).Infoln("body is written")
p.Write([]byte("\n"))
p.Close()
c.Wait()
w.Del(true)
w.Close()
}





/*:197*/


//line amail.w:2729

}
}
w.UnreadEvent(ev)
}
}(msg)



/*:193*/


//line amail.w:2711




/*:192*/


//line amail.w:1064

name:=fmt.Sprintf("Amail/%s/New",box.name)


/*53:*/


//line amail.w:748

glog.V(debug).Infoln("printing a name for a window")
if err:=w.WriteCtl("name %s",strings.Replace(name," ","␣",-1));err!=nil{
glog.Errorf("can't write to 'ctl' file: %s\n",err)
}



/*:53*/


//line amail.w:1066

continue
case"Search":
glog.V(debug).Infof("search argument: '%s'\n",ev.Arg)


/*147:*/


//line amail.w:1931

{
msgs:=box.search(ev.Arg)


/*143:*/


//line amail.w:1843

glog.V(debug).Infof("clean the '%s' mailbox's window\n",box.name)
clean(box.w)



/*:143*/


//line amail.w:1934



/*129:*/


//line amail.w:1649

box.pos= 0
ontop= false



/*:129*/



/*135:*/


//line amail.w:1744

{
glog.V(debug).Infof("clean window-specific stuff of the '%s' mailbox\n",box.name)
close(box.irfch)
box.irfch= make(chan*struct{seek bool;msgs messages},100)
pcount= 0
ontop= false
}



/*:135*/


//line amail.w:1935

name:=fmt.Sprintf("Amail/%s/Search(%s)",box.name,strings.Replace(ev.Arg," ","␣",-1))
w:=box.w
box.thread= false
box.shownew= true


/*53:*/


//line amail.w:748

glog.V(debug).Infoln("printing a name for a window")
if err:=w.WriteCtl("name %s",strings.Replace(name," ","␣",-1));err!=nil{
glog.Errorf("can't write to 'ctl' file: %s\n",err)
}



/*:53*/


//line amail.w:1940

glog.V(debug).Infof("len of msgs: %v\n",len(msgs))
box.rfch<-&struct{seek bool;msgs messages}{false,msgs}
}



/*:147*/


//line amail.w:1070

continue
default:
box.w.UnreadEvent(ev)
continue
}


/*141:*/


//line amail.w:1809

name:="Amail/"+box.name
w:=box.w


/*53:*/


//line amail.w:748

glog.V(debug).Infoln("printing a name for a window")
if err:=w.WriteCtl("name %s",strings.Replace(name," ","␣",-1));err!=nil{
glog.Errorf("can't write to 'ctl' file: %s\n",err)
}



/*:53*/


//line amail.w:1812




/*:141*/


//line amail.w:1076



/*142:*/


//line amail.w:1815

glog.V(debug).Infof("write a tag of the '%s' mailbox's window\n",box.name)
if err:=writeTag(box.w,fmt.Sprintf(" %sMail Delmesg %s%s %s Search ",
func()string{
if box.deleted> 0{
return"Put "
}
return""
}(),
func()string{
if box.thread{
if box.shownew{
return"ShowNew "
}else{
return"ShowAll "
}
}else if box.shownew||!box.showthreads{
return"Thread "
}
return""
}(),
func()string{if box.shownew{return"ShowAll"}else{return"ShowNew"}}(),
func()string{if box.showthreads{return"ShowPlain"}else{return"ShowThreads"}}()))
err!=nil{
glog.Errorf("can't set a tag of the '%s' box's window: %v\n",box.name,err)
}



/*:142*/


//line amail.w:1077



/*143:*/


//line amail.w:1843

glog.V(debug).Infof("clean the '%s' mailbox's window\n",box.name)
clean(box.w)



/*:143*/


//line amail.w:1078



/*50:*/


//line amail.w:718

if w!=nil{
glog.V(debug).Infoln("setting the window to clean state")
if err:=w.WriteCtl("clean");err!=nil{
glog.Errorf("can't write to 'ctl' file: %s\n",err)
}
}



/*:50*/


//line amail.w:1079



/*129:*/


//line amail.w:1649

box.pos= 0
ontop= false



/*:129*/



/*135:*/


//line amail.w:1744

{
glog.V(debug).Infof("clean window-specific stuff of the '%s' mailbox\n",box.name)
close(box.irfch)
box.irfch= make(chan*struct{seek bool;msgs messages},100)
pcount= 0
ontop= false
}



/*:135*/


//line amail.w:1080



/*122:*/


//line amail.w:1574

{
glog.V(debug).Infof("inform the '%s' mailbox to print messages\n",box.name)


/*121:*/


//line amail.w:1564

var src messages
if box.shownew{
src= box.unread
}else{
src= box.all
}



/*:121*/


//line amail.w:1577

msgs:=append(messages{},src...)
if box.threadMode(){
src= msgs
msgs= make(messages,0,len(src))
for len(src)> 0{
msg:=src[0]


/*104:*/


//line amail.w:1356

for msg.parent!=nil{
msg= msg.parent
}



/*:104*/


//line amail.w:1584

glog.V(debug).Infof("root of thread: '%s/%d'\n",msg.box.name,msg.id)


/*123:*/


//line amail.w:1594

msgs= append(msgs,msg)
if p,ok:=src.Search(msg.id);ok&&src[p]==msg{
glog.V(debug).Infof("removing '%d' from src\n",src[p].id)
src.Delete(p)
}
msgs,src= getchildren(msg,msgs,src)




/*:123*/


//line amail.w:1586

}
}
box.rfch<-&struct{seek bool;msgs messages}{false,msgs}
}



/*:122*/


//line amail.w:1081

continue
}else if(ev.Type&goacme.Look)==goacme.Look{


/*84:*/


//line amail.w:1147

msgs:=make(map[string][]int)



/*:84*/


//line amail.w:1084

if(ev.Type&goacme.Tag)==goacme.Tag{
s:=ev.Text


/*82:*/


//line amail.w:1121

{
glog.V(debug).Infof("looking a message number in '%s'\n",s)
s= strings.TrimLeft(s,levelmark+deleted)
f:=strings.Split(s,"/")
glog.V(debug).Infof("parts of message path: '%v'\n",f)
num:=0
for i,v:=range f{
var err error
if num,err= strconv.Atoi(strings.TrimRight(v,newmark));err==nil{
name:=box.name
if i> 0{
name= strings.Join(f[:i],"/")
glog.V(debug).Infof("the message number is '%d' in the '%s' mailbox\n",num,name)
}


/*85:*/


//line amail.w:1151

glog.V(debug).Infof("sending a signal to open a window with the '%d' message of the '%s' mailbox\n",num,name)
msgs[name]= append(msgs[name],num)



/*:85*/


//line amail.w:1136

break
}
}
}



/*:82*/


//line amail.w:1087

}else{


/*81:*/


//line amail.w:1102

glog.V(debug).Infof("event: %v\n",ev)
if err:=box.w.WriteAddr("#%d,#%d",ev.Begin,ev.End);err!=nil{
glog.Errorf("can't write to 'addr': %s\n",err)
}else if xdata,err:=box.w.File("xdata");err!=nil{
glog.Errorf("can't open 'xdata' file: %s\n",err)
}else{
b:=bufio.NewReader(xdata)
for s,err:=b.ReadString('\n');err==nil||err==io.EOF;s,err= b.ReadString('\n'){


/*82:*/


//line amail.w:1121

{
glog.V(debug).Infof("looking a message number in '%s'\n",s)
s= strings.TrimLeft(s,levelmark+deleted)
f:=strings.Split(s,"/")
glog.V(debug).Infof("parts of message path: '%v'\n",f)
num:=0
for i,v:=range f{
var err error
if num,err= strconv.Atoi(strings.TrimRight(v,newmark));err==nil{
name:=box.name
if i> 0{
name= strings.Join(f[:i],"/")
glog.V(debug).Infof("the message number is '%d' in the '%s' mailbox\n",num,name)
}


/*85:*/


//line amail.w:1151

glog.V(debug).Infof("sending a signal to open a window with the '%d' message of the '%s' mailbox\n",num,name)
msgs[name]= append(msgs[name],num)



/*:85*/


//line amail.w:1136

break
}
}
}



/*:82*/


//line amail.w:1111

if err==io.EOF{
break
}
}

}



/*:81*/


//line amail.w:1089

}
if len(msgs)!=0{


/*87:*/


//line amail.w:1167

lch<-&msgs



/*:87*/


//line amail.w:1092

continue
}
}
box.w.UnreadEvent(ev)



/*:80*/



/*120:*/


//line amail.w:1546

case v:=<-box.rfch:
box.irfch<-v

case v:=<-box.irfch:
glog.V(debug).Infof("a signal to print message of the '%s' mailbox window has been received\n",box.name)
if box.w==nil{
glog.V(debug).Infof("a window of the '%s' mailbox doesn't exist, ignore the signal\n",box.name)
continue
}
if box.threadMode()&&!counted{
glog.V(debug).Infof("counting of threads of the '%s' mailbox is not finished, ignore the signal\n",box.name)
continue
}


/*132:*/


//line amail.w:1679

{
glog.V(debug).Infof("printing of messages of the '%s' mailbox from v.msgs, len(v.msgs): %d, with seeking a position: %v\n",box.name,len(v.msgs),v.seek)
f,err:=box.w.File("data")
if err!=nil{
glog.Errorf("can't open 'data' file of the '%s' messgebox: %v\n",box.name,err)
continue
}
if v.seek{


/*142:*/


//line amail.w:1815

glog.V(debug).Infof("write a tag of the '%s' mailbox's window\n",box.name)
if err:=writeTag(box.w,fmt.Sprintf(" %sMail Delmesg %s%s %s Search ",
func()string{
if box.deleted> 0{
return"Put "
}
return""
}(),
func()string{
if box.thread{
if box.shownew{
return"ShowNew "
}else{
return"ShowAll "
}
}else if box.shownew||!box.showthreads{
return"Thread "
}
return""
}(),
func()string{if box.shownew{return"ShowAll"}else{return"ShowNew"}}(),
func()string{if box.showthreads{return"ShowPlain"}else{return"ShowThreads"}}()))
err!=nil{
glog.Errorf("can't set a tag of the '%s' box's window: %v\n",box.name,err)
}



/*:142*/


//line amail.w:1688

msg:=v.msgs[0]


/*145:*/


//line amail.w:1861



/*121:*/


//line amail.w:1564

var src messages
if box.shownew{
src= box.unread
}else{
src= box.all
}



/*:121*/


//line amail.w:1862

addr:=fmt.Sprintf("0/^[%s]*(%s)?%s%d(%s)?\\/.*\\n\t.*\\n/",
escape(levelmark),
escape(deleted),
func()string{if box.name!=msg.box.name{return msg.box.name+"/"};return""}(),
msg.id,
escape(newmark))
glog.V(debug).Infoln("refreshed message addr:",addr)
if err:=box.w.WriteAddr(addr);err!=nil{
glog.V(debug).Infof("the '%d' message is not found in the window\n",msg.id)
if box.threadMode(){


/*146:*/


//line amail.w:1897

if msg.parent!=nil{
glog.V(debug).Infof("msg '%d' has a parent, looking for a position to print\n",msg.id)
m:=msg
msg:=m.parent
found:=false
for!found{


/*115:*/


//line amail.w:1506

ch:=make(chan idmessages)
glog.V(debug).Infof("getting children for '%s'\n",msg.messageid)
idch<-struct{id string;val interface{}}{msg.messageid,ch}
children:=<-ch



/*:115*/


//line amail.w:1904

if children==nil{
break
}
for i,v:=range children{
if v==m{
if i==0{
found= true
}
break
}
msg= v
}
}
if err:=box.w.WriteAddr("0/^[%s]*%s%s%d(%s)?\\/.*\\n\t.*\\n/+#0",
escape(levelmark),
func()string{if msg.deleted{return escape(deleted)};return""}(),
func()string{if box.name!=msg.box.name{return msg.box.name+"/"};return""}(),
msg.id,
escape(newmark));err!=nil{
glog.V(debug).Infof("can't write to 'addr': %s\n",err)
}
}else if err:=box.w.WriteAddr("#0-");err!=nil{
glog.Errorf("can't write to 'addr' file: %s\n",err)
}



/*:146*/


//line amail.w:1873

}else if p,ok:=src.Search(msg.id);!ok{
glog.V(debug).Infof("the '%d' message is not found\n",msg.id)
}else if p==0{
if err:=box.w.WriteAddr("#0-");err!=nil{
glog.Errorf("can't write to 'addr' file: %s\n",err)
}
}else if p==len(src)-1{
if err:=box.w.WriteAddr("$");err!=nil{
glog.Errorf("can't write to 'addr' file: %s\n",err)
}
}else if err:=box.w.WriteAddr("0/^[%s]*%s%s%d(%s)?\\/.*\\n\t.*\\n/+#0",
escape(levelmark),
func()string{if src[p-1].deleted{return escape(deleted)};return""}(),
func()string{if box.name!=src[p-1].box.name{return src[p-1].box.name+"/"};return""}(),
src[p-1].id,
escape(newmark));err!=nil{
glog.V(debug).Infof("can't write to 'addr': %s\n",err)
}
}



/*:145*/


//line amail.w:1690

}else if err:=box.w.WriteAddr("$");err!=nil{
glog.Errorf("can't write to 'addr' file: %s\n",err)
continue
}
w:=box.w
glog.V(debug).Infof("printing of messages of the '%s' mailbox\n",box.name)
buf:=make([]byte,0,0x8000)


/*134:*/


//line amail.w:1726

c:=0
for len(v.msgs)> 0&&c<100{
msg:=v.msgs[0]
glog.V(debug).Infof("printing of message with id: %v\n",msg.id)
if box.threadMode(){


/*136:*/


//line amail.w:1754

{
for p:=msg.parent;p!=nil;p= p.parent{
buf= append(buf,levelmark...)
}
}



/*:136*/


//line amail.w:1732

}
c++


/*90:*/


//line amail.w:1198

glog.V(debug).Infof("compose a header of the '%d' message of the '%s' mailbox\n",msg.id,box.name)
buf= append(buf,fmt.Sprintf("%s%s%d%s/\t%s\t%s\n\t%s\n",
func()string{if msg.deleted{return deleted};return""}(),
func()string{if msg.box.name!=box.name{return fmt.Sprintf("%s/",msg.box.name)};return""}(),
msg.id,
func()string{if msg.unread{return newmark};return""}(),
msg.from,
msg.date,
msg.subject)...)



/*:90*/


//line amail.w:1735

v.msgs= v.msgs[1:]
if v.seek{
break
}
}
pcount+= c



/*:134*/


//line amail.w:1698

if _,err:=f.Write(buf);err!=nil{
glog.Errorf("can't write to 'data' file of the '%s' messgebox: %v\n",box.name,err)
}


/*133:*/


//line amail.w:1712

if!ontop{
glog.V(debug).Infof("pcount: %v, ontop: %v\n",pcount,ontop)


/*93:*/


//line amail.w:1234

glog.V(debug).Infoln("Go to top of the window")
w.WriteAddr("#0")
w.WriteCtl("dot=addr\nshow")



/*:93*/


//line amail.w:1715

if pcount>=100{
ontop= true
}
}




/*:133*/


//line amail.w:1702

if len(v.msgs)> 0{
box.rfch<-&struct{seek bool;msgs messages}{v.seek,v.msgs}
}else{


/*155:*/


//line amail.w:2036

glog.V(debug).Infof("box.deleted:%d\n",box.deleted)
if box.deleted==0{
w:=box.w


/*142:*/


//line amail.w:1815

glog.V(debug).Infof("write a tag of the '%s' mailbox's window\n",box.name)
if err:=writeTag(box.w,fmt.Sprintf(" %sMail Delmesg %s%s %s Search ",
func()string{
if box.deleted> 0{
return"Put "
}
return""
}(),
func()string{
if box.thread{
if box.shownew{
return"ShowNew "
}else{
return"ShowAll "
}
}else if box.shownew||!box.showthreads{
return"Thread "
}
return""
}(),
func()string{if box.shownew{return"ShowAll"}else{return"ShowNew"}}(),
func()string{if box.showthreads{return"ShowPlain"}else{return"ShowThreads"}}()))
err!=nil{
glog.Errorf("can't set a tag of the '%s' box's window: %v\n",box.name,err)
}



/*:142*/


//line amail.w:2040



/*50:*/


//line amail.w:718

if w!=nil{
glog.V(debug).Infoln("setting the window to clean state")
if err:=w.WriteCtl("clean");err!=nil{
glog.Errorf("can't write to 'ctl' file: %s\n",err)
}
}



/*:50*/


//line amail.w:2041

}else{
w:=box.w


/*142:*/


//line amail.w:1815

glog.V(debug).Infof("write a tag of the '%s' mailbox's window\n",box.name)
if err:=writeTag(box.w,fmt.Sprintf(" %sMail Delmesg %s%s %s Search ",
func()string{
if box.deleted> 0{
return"Put "
}
return""
}(),
func()string{
if box.thread{
if box.shownew{
return"ShowNew "
}else{
return"ShowAll "
}
}else if box.shownew||!box.showthreads{
return"Thread "
}
return""
}(),
func()string{if box.shownew{return"ShowAll"}else{return"ShowNew"}}(),
func()string{if box.showthreads{return"ShowPlain"}else{return"ShowThreads"}}()))
err!=nil{
glog.Errorf("can't set a tag of the '%s' box's window: %v\n",box.name,err)
}



/*:142*/


//line amail.w:2044



/*51:*/


//line amail.w:727

if w!=nil{
glog.V(debug).Infoln("setting the window to dirty state")
if err:=w.WriteCtl("dirty");err!=nil{
glog.Errorf("can't write to 'ctl' file: %s\n",err)
}
}



/*:51*/


//line amail.w:2045

}



/*:155*/


//line amail.w:1706

}
}



/*:132*/


//line amail.w:1560





/*:120*/



/*152:*/


//line amail.w:1995

case ids:=<-box.lch:
for _,id:=range ids{
glog.V(debug).Infof("opening a window with the '%d' message of the '%s' mailbox\n",id,box.name)
p,ok:=box.all.Search(id)
if!ok{
glog.V(debug).Infof("the '%d' message of the '%s' mailbox has not found\n",id,box.name)
continue
}
msg:=box.all[p]
if msg.w==nil{
if err:=msg.open();err!=nil{
continue
}
if msg.unread{


/*153:*/


//line amail.w:2021

msg.unread= false
box.unread.DeleteById(id)



/*:153*/


//line amail.w:2010



/*154:*/


//line amail.w:2026

if box.thread||!box.shownew{


/*156:*/


//line amail.w:2049

if box.w!=nil{
glog.V(debug).Infof("removing the '%s' newmark from the '%d' message\n",newmark,id)
addr:=fmt.Sprintf("-/^[%s]*(%s)?%d/",escape(levelmark),escape(deleted),id)
if err:=box.w.WriteAddr(addr);err!=nil{
glog.Errorf("can't write '%s' to 'addr': %s\n",addr,err)
}else if err:=box.w.WriteAddr("/%s/",escape(newmark));err!=nil{
glog.Errorf("can't write to 'addr': %s\n",err)
}else if data,err:=box.w.File("data");err!=nil{
glog.Errorf("can't open 'data' file of the box '%s': %s\n",box.name,err)
}else if _,err:=data.Write([]byte{});err!=nil{
glog.Errorf("can't write to 'data' file of the box '%s': %s\n",box.name,err)
}
}



/*:156*/


//line amail.w:2028

}else{


/*157:*/


//line amail.w:2065

box.removeMessage(id)



/*:157*/


//line amail.w:2030

}


/*155:*/


//line amail.w:2036

glog.V(debug).Infof("box.deleted:%d\n",box.deleted)
if box.deleted==0{
w:=box.w


/*142:*/


//line amail.w:1815

glog.V(debug).Infof("write a tag of the '%s' mailbox's window\n",box.name)
if err:=writeTag(box.w,fmt.Sprintf(" %sMail Delmesg %s%s %s Search ",
func()string{
if box.deleted> 0{
return"Put "
}
return""
}(),
func()string{
if box.thread{
if box.shownew{
return"ShowNew "
}else{
return"ShowAll "
}
}else if box.shownew||!box.showthreads{
return"Thread "
}
return""
}(),
func()string{if box.shownew{return"ShowAll"}else{return"ShowNew"}}(),
func()string{if box.showthreads{return"ShowPlain"}else{return"ShowThreads"}}()))
err!=nil{
glog.Errorf("can't set a tag of the '%s' box's window: %v\n",box.name,err)
}



/*:142*/


//line amail.w:2040



/*50:*/


//line amail.w:718

if w!=nil{
glog.V(debug).Infoln("setting the window to clean state")
if err:=w.WriteCtl("clean");err!=nil{
glog.Errorf("can't write to 'ctl' file: %s\n",err)
}
}



/*:50*/


//line amail.w:2041

}else{
w:=box.w


/*142:*/


//line amail.w:1815

glog.V(debug).Infof("write a tag of the '%s' mailbox's window\n",box.name)
if err:=writeTag(box.w,fmt.Sprintf(" %sMail Delmesg %s%s %s Search ",
func()string{
if box.deleted> 0{
return"Put "
}
return""
}(),
func()string{
if box.thread{
if box.shownew{
return"ShowNew "
}else{
return"ShowAll "
}
}else if box.shownew||!box.showthreads{
return"Thread "
}
return""
}(),
func()string{if box.shownew{return"ShowAll"}else{return"ShowNew"}}(),
func()string{if box.showthreads{return"ShowPlain"}else{return"ShowThreads"}}()))
err!=nil{
glog.Errorf("can't set a tag of the '%s' box's window: %v\n",box.name,err)
}



/*:142*/


//line amail.w:2044



/*51:*/


//line amail.w:727

if w!=nil{
glog.V(debug).Infoln("setting the window to dirty state")
if err:=w.WriteCtl("dirty");err!=nil{
glog.Errorf("can't write to 'ctl' file: %s\n",err)
}
}



/*:51*/


//line amail.w:2045

}



/*:155*/


//line amail.w:2032





/*:154*/


//line amail.w:2011



/*72:*/


//line amail.w:932

glog.V(debug).Infof("sending a snapshot of the '%s' mailbox to refresh the main window\n",box.name)
b:=*box
rfch<-&b



/*:72*/


//line amail.w:2012

}
}else{
glog.V(debug).Infof("a window of the '%d' message of the '%s' already exists, just show it\n",id,box.name)
msg.w.WriteCtl("dot=addr\nshow")
}
}



/*:152*/


//line amail.w:856

}
}
}



/*:64*/



/*70:*/


//line amail.w:907

func(box*mailbox)deleteMessage(id int){
glog.V(debug).Infof("deleting the '%d' message from the '%s' mailbox\n",id,box.name)
box.unread.DeleteById(id)
if msg,ok:=box.all.DeleteById(id);ok{
box.total--
if msg.deleted{
box.deleted--
}


/*115:*/


//line amail.w:1506

ch:=make(chan idmessages)
glog.V(debug).Infof("getting children for '%s'\n",msg.messageid)
idch<-struct{id string;val interface{}}{msg.messageid,ch}
children:=<-ch



/*:115*/


//line amail.w:916



/*116:*/


//line amail.w:1513

glog.V(debug).Infof("cleaning up the '%d' message\n",msg.id)
if msg!=nil{
idch<-struct{id string;val interface{}}{id:msg.messageid}
}



/*:116*/


//line amail.w:917



/*157:*/


//line amail.w:2065

box.removeMessage(id)



/*:157*/


//line amail.w:918

if box.threadMode(){


/*137:*/


//line amail.w:1763

{
if len(children)!=0{
var msgs messages
var src messages
for _,msg:=range children{


/*123:*/


//line amail.w:1594

msgs= append(msgs,msg)
if p,ok:=src.Search(msg.id);ok&&src[p]==msg{
glog.V(debug).Infof("removing '%d' from src\n",src[p].id)
src.Delete(p)
}
msgs,src= getchildren(msg,msgs,src)




/*:123*/


//line amail.w:1769

}
for _,msg:=range msgs{


/*126:*/


//line amail.w:1628

{
glog.V(debug).Infof("inform the '%s' mailbox to print a message '%d'\n",box.name,msg.id)
box.rfch<-&struct{seek bool;msgs messages}{true,append(messages{},msg)}
}



/*:126*/


//line amail.w:1772

}
}
}



/*:137*/


//line amail.w:920

}
}
}



/*:70*/



/*71:*/


//line amail.w:926

func(box*mailbox)threadMode()bool{
return box.thread||box.showthreads&&!box.shownew
}



/*:71*/



/*106:*/


//line amail.w:1373

func(this idmessages)Search(messageid string)(int,bool){
pos:=sort.Search(len(this),func(i int)bool{return this[i].messageid<=messageid});
if pos!=len(this)&&this[pos].messageid==messageid{
return pos,true
}
return pos,false
}



/*:106*/



/*107:*/


//line amail.w:1383

func(this*idmessages)Insert(msg*message,pos int){
*this= append(*this,nil)
copy((*this)[pos+1:],(*this)[pos:])
(*this)[pos]= msg
}



/*:107*/



/*108:*/


//line amail.w:1393

func(this*idmessages)SearchInsert(msg*message)(int,bool){
pos,ok:=this.Search(msg.messageid)
if ok{
return pos,false
}
this.Insert(msg,pos)
return pos,true
}



/*:108*/



/*110:*/


//line amail.w:1407

func(this*idmessages)Delete(pos int)(*message,bool){
if pos<0||pos> len(*this)-1{
return nil,false
}
msg:=(*this)[pos]
*this= append((*this)[:pos],(*this)[pos+1:]...)
return msg,true
}



/*:110*/



/*113:*/


//line amail.w:1474

func(this idmessages)Len()int{
return len(this)
}

func(this idmessages)Less(i,j int)bool{
return this[i].date.Unix()<this[j].date.Unix()
}

func(this idmessages)Swap(i,j int){
t:=this[i]
this[i]= this[j]
this[j]= t
}




/*:113*/



/*124:*/


//line amail.w:1605

func getchildren(msg*message,dst messages,src messages)(messages,messages){


/*115:*/


//line amail.w:1506

ch:=make(chan idmessages)
glog.V(debug).Infof("getting children for '%s'\n",msg.messageid)
idch<-struct{id string;val interface{}}{msg.messageid,ch}
children:=<-ch



/*:115*/


//line amail.w:1607

for _,v:=range children{
dst= append(dst,v)
if p,ok:=src.Search(v.id);ok&&src[p]==v{
glog.V(debug).Infof("removing '%d' from src\n",src[p].id)
src.Delete(p)
}
dst,src= getchildren(v,dst,src)
}
return dst,src
}



/*:124*/



/*139:*/


//line amail.w:1788

func writeTag(w*goacme.Window,t string)error{
if w==nil{
return nil
}
tag,err:=w.File("tag")
if err!=nil{
return err
}
if err:=w.WriteCtl("cleartag");err!=nil{
return err
}
_,err= tag.Write([]byte(t))
return err
}



/*:139*/



/*144:*/


//line amail.w:1848

func clean(w*goacme.Window){
if err:=w.WriteAddr("0,$");err!=nil{
glog.Errorf("can't write to 'addr' file: %s\n",err)
}else if data,err:=w.File("data");err!=nil{
glog.Errorf("can't open 'data' file: %s\n",err)
}else if _,err:=data.Write([]byte(""));err!=nil{
glog.Errorf("can't write to 'data' file: %s\n",err)
}
}



/*:144*/



/*148:*/


//line amail.w:1946

func(box*mailbox)search(str string)(msgs messages){
if len(str)==0{
return
}
f,err:=box.fid.Walk("search")
if err==nil{
err= f.Open(plan9.ORDWR)
}
if err!=nil{
glog.Errorf("can't open 'search' file: %s\n",err)
return
}
defer f.Close()
if _,err:=f.Write([]byte(str));err!=nil{
glog.Errorf("can't write to 'search' file: %s\n",err)
}
b:=bufio.NewReader(f)
for s,err:=b.ReadString(' ');err==nil||err==io.EOF;s,err= b.ReadString(' '){
s= strings.TrimSpace(s)
glog.V(debug).Infoln("search: ",s)
if num,err:=strconv.Atoi(s);err==nil{
if p,ok:=box.all.Search(num);ok{
msgs.Insert(box.all[p],0)
}
}
if err==io.EOF{
break
}
}
return
}



/*:148*/



/*158:*/


//line amail.w:2069

func(box*mailbox)removeMessage(id int){
if box.w==nil{
return
}
glog.V(debug).Infof("removing the '%d' message\n",id)
addr:=fmt.Sprintf("-/^[%s]*(%s)?%d(%s)?\\//",escape(levelmark),escape(deleted),id,escape(newmark))
if err:=box.w.WriteAddr(addr);err!=nil{
glog.V(debug).Infof("can't write '%s' to 'addr': %s\n",addr,err)
}else if err:=box.w.WriteAddr(".,./.*\\n(\t.*\\n)*/");err!=nil{
glog.Errorf("can't write to 'addr': %s\n",err)
}else if data,err:=box.w.File("data");err!=nil{
glog.Errorf("can't open 'data' file of the box '%s': %s\n",box.name,err)
}else if _,err:=data.Write([]byte{});err!=nil{
glog.Errorf("can't write to 'data' file of the box '%s': %s\n",box.name,err)
}
}



/*:158*/



/*161:*/


//line amail.w:2106

func split(s string)(strs[]string){
f:=strings.Fields(s)
m:=""
for _,v:=range f{
if strings.Contains(v,"@"){
m+= v
strs= append(strs,m)
m= ""
}else if v!="''"{
m+= v+" "
}
}
return
}



/*:161*/



/*162:*/


//line amail.w:2123

func(this*message)open()(err error){
glog.V(debug).Infof("open: trying to open '%d' directory\n",this.id)
bfid,err:=this.box.fid.Walk(fmt.Sprintf("%d",this.id))
if err!=nil{
glog.Errorf("can't walk to '%s/%d': %v\n",this.box.name,this.id,err)
return err
}
defer bfid.Close()
if this.w==nil{
if this.w,err= goacme.New();err!=nil{
glog.Errorf("can't create a window: %v\n",err)
return err
}
msg:=this


/*166:*/


//line amail.w:2191

go func(){
glog.V(debug).Infof("starting a goroutine to process events from the '%d' message's window\n",this.id)
for ev,err:=this.w.ReadEvent();err==nil;ev,err= this.w.ReadEvent(){
if ev.Origin!=goacme.Mouse{
this.w.UnreadEvent(ev)
continue
}
quote:=false
replyall:=false
if(ev.Type&goacme.Execute)==goacme.Execute{
switch ev.Text{
case"Del":
this.w.UnreadEvent(ev)
this.w.Close()
this.w= nil
return
case"Delmesg":
if!this.deleted{
this.deleted= true
this.box.deleted++
this.w.Del(true)
this.w.Close()
this.w= nil
box:=this.box
msg:=this


/*126:*/


//line amail.w:1628

{
glog.V(debug).Infof("inform the '%s' mailbox to print a message '%d'\n",box.name,msg.id)
box.rfch<-&struct{seek bool;msgs messages}{true,append(messages{},msg)}
}



/*:126*/


//line amail.w:2217

return
}
continue
case"UnDelmesg":
if this.deleted{
this.deleted= false
this.box.deleted--


/*163:*/


//line amail.w:2156

this.writeTag()



/*:163*/


//line amail.w:2225

box:=this.box
msg:=this


/*126:*/


//line amail.w:1628

{
glog.V(debug).Infof("inform the '%s' mailbox to print a message '%d'\n",box.name,msg.id)
box.rfch<-&struct{seek bool;msgs messages}{true,append(messages{},msg)}
}



/*:126*/


//line amail.w:2228

}
continue
case"Text":
if len(this.text)!=0&&len(this.html)!=0{
this.showhtml= false
this.open()
}
continue
case"Html":
if len(this.text)!=0&&len(this.html)!=0{
this.showhtml= true
this.open()
}
continue
case"Browser":


/*180:*/


//line amail.w:2506

{


/*185:*/


//line amail.w:2612



/*189:*/


//line amail.w:2633

once.Do(func(){

/*177:*/


//line amail.w:2450

if home= os.Getenv("home");len(home)==0{
if home= os.Getenv("HOME");len(home)==0{
glog.Errorln("can't get a home directory from the environment, the home is assumed '/'")
home= "/"
}
}




/*:177*/



/*186:*/


//line amail.w:2616

if u,err:=user.Current();err!=nil{
glog.Errorf("can't get a name of the current user: %v\n",err)
}else{
cuser= u.Username
}




/*:186*/



/*200:*/


//line amail.w:2922

if plan9dir= os.Getenv("PLAN9");len(plan9dir)==0{
glog.Errorln("can't get PLAN9 directory from the environment, the plan9dir is assumed '/usr/local/plan9'")
plan9dir= "/usr/local/plan9"
}



/*:200*/


//line amail.w:2634
})



/*:189*/


//line amail.w:2613




/*:185*/


//line amail.w:2508

dir:=fmt.Sprintf("%s/amail-%s/%s/%d",os.TempDir(),cuser,this.box.name,this.id)
if err:=os.MkdirAll(dir,0700);err!=nil{
glog.Errorf("can't create a directory '%s': %v\n",dir,err)
continue
}

if len(this.files)==0{
if err:=saveFile(fmt.Sprintf("%s/%s/%d/%s",srv,this.box.name,this.id,this.html),
fmt.Sprintf("%s/%d.html",dir,this.id));err!=nil{
continue
}
}else{
if err:=this.fixFile(dir);err!=nil{
continue
}
for _,v:=range this.files{
saveFile(fmt.Sprintf("%s/%s/%d/%s/body",srv,this.box.name,this.id,v.path),
fmt.Sprintf("%s/%s",dir,v.name))
}

}

if p,err:=goplumb.Open("send",plan9.OWRITE);err!=nil{
glog.Errorf("can't open plumbing port 'send': %v\n",err)
}else if err:=p.SendText("amail","web",dir,fmt.Sprintf("file://%s/%d.html",dir,this.id));err!=nil{
glog.Errorf("can't plumb a message '%s': %v\n",fmt.Sprintf("file://%s/%d.html",dir,this.id),err)
}
}



/*:180*/


//line amail.w:2244

continue
case"Save":


/*190:*/


//line amail.w:2637

{
if len(ev.Arg)==0{
continue
}
f,err:=this.box.fid.Walk("ctl")
if err==nil{
err= f.Open(plan9.OWRITE)
}
if err!=nil{
glog.Errorf("can't open 'ctl': %v\n",err)
continue
}
bs:=strings.Fields(ev.Arg)
for _,v:=range bs{
s:=fmt.Sprintf("save %s %d/",v,this.id)
if _,err:=f.Write([]byte(s));err!=nil{
glog.Errorf("can't write '%s' to 'ctl': %v\n",s,err)
}
}
f.Close()
}




/*:190*/


//line amail.w:2247

continue
case"Q":
quote= true
fallthrough
case"Reply","Replyall":
if ev.Text=="Reply"{
args:=strings.Fields(ev.Arg)
for _,v:=range args{
if v=="all"{
replyall= true
}
}
}else if ev.Text=="Replyall"{
replyall= true
}


/*191:*/


//line amail.w:2662

{


/*192:*/


//line amail.w:2702

w,err:=goacme.New()
if err!=nil{
glog.Errorf("can't create a window: %v\n",err)
continue
}
if err:=writeTag(w," Look Post Undo ");err!=nil{
glog.Errorf("can't write a tag for a new message window: %v\n",err)
}


/*193:*/


//line amail.w:2714

go func(msg*message){
glog.V(debug).Infoln("starting a goroutine to process events from a composed mail window")
for ev,err:=w.ReadEvent();err==nil;ev,err= w.ReadEvent(){
if ev.Origin!=goacme.Mouse{
w.UnreadEvent(ev)
continue
}
if(ev.Type&goacme.Execute)==goacme.Execute{
switch ev.Text{
case"Del":
w.UnreadEvent(ev)
w.Close()
return
case"Post":


/*197:*/


//line amail.w:2801

{


/*199:*/


//line amail.w:2918



/*189:*/


//line amail.w:2633

once.Do(func(){

/*177:*/


//line amail.w:2450

if home= os.Getenv("home");len(home)==0{
if home= os.Getenv("HOME");len(home)==0{
glog.Errorln("can't get a home directory from the environment, the home is assumed '/'")
home= "/"
}
}




/*:177*/



/*186:*/


//line amail.w:2616

if u,err:=user.Current();err!=nil{
glog.Errorf("can't get a name of the current user: %v\n",err)
}else{
cuser= u.Username
}




/*:186*/



/*200:*/


//line amail.w:2922

if plan9dir= os.Getenv("PLAN9");len(plan9dir)==0{
glog.Errorln("can't get PLAN9 directory from the environment, the plan9dir is assumed '/usr/local/plan9'")
plan9dir= "/usr/local/plan9"
}



/*:200*/


//line amail.w:2634
})



/*:189*/


//line amail.w:2919




/*:199*/


//line amail.w:2803

w.Seek(0,0)
w.WriteAddr("0,$")
ff,_:=w.File("xdata")
b:=bufio.NewReader(ff)
var to,cc,bcc,attach,include[]string
var subject string
for{
s,err:=b.ReadString('\n')
if err!=nil{
break
}
s= strings.TrimSpace(s)
if len(s)==0{
// an empty line, the rest is a body of the message
break
}
p:=strings.Index(s,":")
if p!=-1{
f:=strings.Split(s[p+1:],",")
for i,_:=range f{
f[i]= strings.TrimSpace(f[i])
}
switch strings.ToLower(s[:p]){
case"from","to":
to= append(to,f...)
case"cc":
cc= append(cc,f...)
case"bcc":
bcc= append(bcc,f...)
case"attach":
attach= append(attach,f...)
case"include":
include= append(include,f...)
case"subject":
subject= fmt.Sprintf("%q",strings.TrimSpace(s[p+1:]))
}
}else{
// recipient addresses can be written without "to:"
f:=strings.Split(s,",")
for i,_:=range f{
f[i]= strings.TrimSpace(f[i])
}
to= append(to,f...)
}
}
args:=append([]string{},"-8")
if msg!=nil{
args= append(args,"-R",fmt.Sprintf("%s/%d",msg.box.name,msg.id))
}
if len(subject)!=0{
args= append(args,"-s",subject)
}
for _,v:=range include{
args= append(args,"-A",v)
}
for _,v:=range attach{
args= append(args,"-a",v)
}

c:=exec.Command(plan9dir+"/bin/upas/marshal",args...)
p,err:=c.StdinPipe()
if err!=nil{
glog.Errorf("can't get a stdin pipe: %v\n",err)
continue
}
if err:=c.Start();err!=nil{
glog.Errorf("can't start 'upas/marshal': %v\n",err)
continue
}
if len(to)!=0{
if _,err:=fmt.Fprintln(p,"To:",strings.Join(to,","));err!=nil{
glog.Errorf("can't write 'to' fields to 'upas/marshal': %v\n",err)
continue
}
}
glog.V(debug).Infoln("to is written")
if len(cc)!=0{
if _,err:=fmt.Fprintln(p,"CC:",strings.Join(cc,","));err!=nil{
glog.Errorf("can't write 'cc' fields to 'upas/marshal': %v\n",err)
continue
}
}
glog.V(debug).Infoln("cc is written")
if len(bcc)!=0{
if _,err:=fmt.Fprintln(p,"BCC:",strings.Join(bcc,","));err!=nil{
glog.Errorf("can't write 'bcc' fields to  'upas/marshal': %v\n",err)
continue
}
}
glog.V(debug).Infoln("bcc is written")
for s,err:=b.ReadString('\n');err==nil||err==io.EOF;s,err= b.ReadString('\n'){
glog.V(debug).Infof("writing '%s':%v",s,err)

p.Write([]byte(s))
if err==io.EOF{
p.Write([]byte("\n"))
break
}
}
glog.V(debug).Infoln("body is written")
p.Write([]byte("\n"))
p.Close()
c.Wait()
w.Del(true)
w.Close()
}





/*:197*/


//line amail.w:2729

}
}
w.UnreadEvent(ev)
}
}(msg)



/*:193*/


//line amail.w:2711




/*:192*/


//line amail.w:2664

name:=fmt.Sprintf("Amail/%s/%d/%sReply%s",
this.box.name,
this.id,
func()string{if quote{return"Q"};return""}(),
func()string{if replyall{return"all"};return""}())


/*53:*/


//line amail.w:748

glog.V(debug).Infoln("printing a name for a window")
if err:=w.WriteCtl("name %s",strings.Replace(name," ","␣",-1));err!=nil{
glog.Errorf("can't write to 'ctl' file: %s\n",err)
}



/*:53*/


//line amail.w:2670

buf:=make([]byte,0,0x8000)
buf= append(buf,fmt.Sprintf("To: %s\n",this.from)...)
if replyall{
for _,v:=range this.to{
buf= append(buf,fmt.Sprintf("To: %s\n",v)...)
}
for _,v:=range this.cc{
buf= append(buf,fmt.Sprintf("To: %s\n",v)...)
}
}
buf= append(buf,fmt.Sprintf("Subject: %s%s\n",
func()string{
if!strings.Contains(this.subject,"Re:"){
return"Re: "
}
return""
}(),
this.subject)...)
if quote{
buf= append(buf,'\n')


/*194:*/


//line amail.w:2737

if len(this.text)!=0{
fn:=fmt.Sprintf("%d/%s",this.id,this.text)
f,err:=this.box.fid.Walk(fn)
if err==nil{
err= f.Open(plan9.OREAD)
}
if err!=nil{
glog.Errorf("can't open '%s/%s/%s': %v\n",srv,this.box.name,fn)
continue
}


/*195:*/


//line amail.w:2755

{
b:=bufio.NewReader(f)
for s,err:=b.ReadString('\n');err==nil||err==io.EOF;s,err= b.ReadString('\n'){
buf= append(buf,'>',' ')
buf= append(buf,s...)
if err==io.EOF{
break
}
}
}



/*:195*/


//line amail.w:2748

f.Close()
}else if len(this.html)!=0{


/*196:*/


//line amail.w:2768

{
c1:=exec.Command("9p","read",fmt.Sprintf("%s/%s/%d/%s",srv,this.box.name,this.id,this.html))
c2:=exec.Command("htmlfmt","-cutf-8")
f,err:=c2.StdoutPipe()
if err!=nil{
glog.Errorf("can't get a stdout pipe: %v\n",err)
}
c2.Stdin,err= c1.StdoutPipe()
if err!=nil{
glog.Errorf("can't get a stdout pipe: %v\n",err)
f.(io.Closer).Close()
continue
}
if err= c2.Start();err!=nil{
glog.Errorf("can't start 'htmlfmt': %v\n",err)
f.(io.Closer).Close()
continue
}
if err= c1.Start();err!=nil{
glog.Errorf("can't run '9p': %v\n",err)
c2.Wait()
f.(io.Closer).Close()
continue
}


/*195:*/


//line amail.w:2755

{
b:=bufio.NewReader(f)
for s,err:=b.ReadString('\n');err==nil||err==io.EOF;s,err= b.ReadString('\n'){
buf= append(buf,'>',' ')
buf= append(buf,s...)
if err==io.EOF{
break
}
}
}



/*:195*/


//line amail.w:2793

c1.Wait()
c2.Wait()
f.(io.Closer).Close()
}



/*:196*/


//line amail.w:2751

}



/*:194*/


//line amail.w:2691

}else{
buf= append(buf,fmt.Sprintf("Include: Mail/%s/%d/raw\n",this.box.name,this.id)...)

}
buf= append(buf,'\n')
w.Write(buf)


/*93:*/


//line amail.w:1234

glog.V(debug).Infoln("Go to top of the window")
w.WriteAddr("#0")
w.WriteCtl("dot=addr\nshow")



/*:93*/


//line amail.w:2698

}



/*:191*/


//line amail.w:2263

continue
}
}else if(ev.Type&goacme.Look)==goacme.Look{
}
this.w.UnreadEvent(ev)

}
}()



/*:166*/


//line amail.w:2138

}else{


/*167:*/


//line amail.w:2274

glog.V(debug).Infof("clean the '%s/%d' message's window\n",this.box.name,this.id)
clean(this.w)



/*:167*/


//line amail.w:2140

}
buf:=make([]byte,0,0x8000)


/*165:*/


//line amail.w:2181

{
glog.V(debug).Infof("composing a header of the '%d' message\n",this.id)
buf= append(buf,fmt.Sprintf("From: %s\nDate: %s\nTo: %s\n%sSubject: %s\n\n\n",
this.from,this.date,strings.Join(this.to,", "),
func()string{if len(this.cc)!=0{return fmt.Sprintf("CC: %s\n",strings.Join(this.cc,", "))};return""}(),
this.subject)...)
}



/*:165*/


//line amail.w:2143



/*172:*/


//line amail.w:2315

{
if len(this.text)==0&&len(this.html)==0{
if err= this.bodyPath(bfid,"");err!=nil{
glog.Errorf("can't ged a body path of '%d': %v\n",this.id,err)
}
glog.V(debug).Infof("paths for bodies of the '%d' message have been found: text-'%s', html-'%s'\n",
this.id,this.text,this.html)

}
if len(this.text)!=0&&!this.showhtml{
glog.V(debug).Infof("using a path for a text body of the '%d' message: '%s'\n",this.id,this.text)
if buf,err= readAll(bfid,this.text,buf);err!=nil{
glog.Errorf("can't read '%s': %v\n",this.text,err)
return
}
}else if len(this.html)!=0{
glog.V(debug).Infof("using a path for a html body of the '%d' message: '%s'\n",this.id,this.html)
this.w.Write(buf)
buf= nil
c1:=exec.Command("9p","read",fmt.Sprintf("%s/%s/%d/%s",srv,this.box.name,this.id,this.html))
c2:=exec.Command("htmlfmt","-cutf-8")
c2.Stdout,_= this.w.File("body")
c2.Stdin,err= c1.StdoutPipe()
if err!=nil{
glog.Errorf("can't get a stdout pipe: %v\n",err)
return
}
if err= c2.Start();err!=nil{
glog.Errorf("can't start 'htmlfmt': %v\n",err)
return
}
if err= c1.Run();err!=nil{
glog.Errorf("can't run '9p': %v\n",err)
c2.Wait()
return
}
if err= c2.Wait();err!=nil{
glog.Errorf("can't wait of completion 'htmlfmt': %v\n",err)
return
}
}


/*176:*/


//line amail.w:2446



/*189:*/


//line amail.w:2633

once.Do(func(){

/*177:*/


//line amail.w:2450

if home= os.Getenv("home");len(home)==0{
if home= os.Getenv("HOME");len(home)==0{
glog.Errorln("can't get a home directory from the environment, the home is assumed '/'")
home= "/"
}
}




/*:177*/



/*186:*/


//line amail.w:2616

if u,err:=user.Current();err!=nil{
glog.Errorf("can't get a name of the current user: %v\n",err)
}else{
cuser= u.Username
}




/*:186*/



/*200:*/


//line amail.w:2922

if plan9dir= os.Getenv("PLAN9");len(plan9dir)==0{
glog.Errorln("can't get PLAN9 directory from the environment, the plan9dir is assumed '/usr/local/plan9'")
plan9dir= "/usr/local/plan9"
}



/*:200*/


//line amail.w:2634
})



/*:189*/


//line amail.w:2447




/*:176*/


//line amail.w:2357

for _,v:=range this.files{
buf= append(buf,fmt.Sprintf("\n===> %s (%s)\n",v.path,v.mimetype)...)
buf= append(buf,fmt.Sprintf("\t9p read %s/%s/%d/%sbody > %s/%s\n",srv,this.box.name,this.id,v.path,home,v.name)...)
}
}



/*:172*/


//line amail.w:2144

w:=this.w
name:=fmt.Sprintf("Amail/%s/%d",this.box.name,this.id)


/*53:*/


//line amail.w:748

glog.V(debug).Infoln("printing a name for a window")
if err:=w.WriteCtl("name %s",strings.Replace(name," ","␣",-1));err!=nil{
glog.Errorf("can't write to 'ctl' file: %s\n",err)
}



/*:53*/


//line amail.w:2147



/*163:*/


//line amail.w:2156

this.writeTag()



/*:163*/


//line amail.w:2148

w.Write(buf)


/*50:*/


//line amail.w:718

if w!=nil{
glog.V(debug).Infoln("setting the window to clean state")
if err:=w.WriteCtl("clean");err!=nil{
glog.Errorf("can't write to 'ctl' file: %s\n",err)
}
}



/*:50*/


//line amail.w:2150



/*93:*/


//line amail.w:1234

glog.V(debug).Infoln("Go to top of the window")
w.WriteAddr("#0")
w.WriteCtl("dot=addr\nshow")



/*:93*/


//line amail.w:2151

return
}



/*:162*/



/*164:*/


//line amail.w:2160

func(this*message)writeTag(){
glog.V(debug).Infof("writing a tag of the '%d' message's window\n",this.id)
if err:=writeTag(this.w,fmt.Sprintf(" Q Reply all %s %s%sSave ",
func()string{if this.deleted{return"UnDelmesg"}else{return"Delmesg"}}(),
func()string{
if len(this.text)==0||len(this.html)==0{
return""
}else if this.showhtml{
return"Text "
}else{
return"Html "
}
}(),
func()string{if len(this.html)!=0{return"Browser "};return""}()))
err!=nil{
glog.Errorf("can't set a tag of the message window: %v",err)
}
}



/*:164*/



/*173:*/


//line amail.w:2365

func(this*message)bodyPath(bfid*client.Fid,path string)error{
glog.V(debug).Infof("getting a path for a body of the '%d' message\n",this.id)
t,err:=readString(bfid,path+"type")
if err!=nil{
return err
}
switch t{
case"message/rfc822",
"text",
"text/plain",
"text/richtext",
"text/tab-separated-values":
if len(this.text)==0{
this.text= path+"body"
glog.V(debug).Infof("a path for a text body of the '%d' message: '%s'\n",this.id,t)
}
case"text/html":
if len(this.html)==0{
this.html= path+"body"
glog.V(debug).Infof("a path for a html body of the '%d' message: '%s'\n",this.id,t)
}
case"multipart/mixed",
"multipart/alternative",
"multipart/related",
"multipart/signed":
for c:=1;;c++{
if err= this.bodyPath(bfid,fmt.Sprintf("%s%d/",path,c));err!=nil{
break
}
}
default:
glog.V(debug).Infof("trying to read '%d/%sfilename'\n",this.id,path)
if n,err:=readString(bfid,path+"filename");err==nil&&len(n)> 0{
f:=&file{name:n,mimetype:t,path:path,}
if cid,ok:=this.getCID(path);ok{
this.cids[cid]= f
}
this.files= append(this.files,f)
}
}
return nil
}



/*:173*/



/*174:*/


//line amail.w:2410

func(this*message)getCID(path string)(string,bool){
src:=fmt.Sprintf("%d/%smimeheader",this.id,path)
glog.V(debug).Infof("getting of cids for path '%s'\n",src)
fid,err:=this.box.fid.Walk(src)
if err==nil{
err= fid.Open(plan9.OREAD)
}
if err!=nil{
glog.Errorf("can't open '%s': %v\n",src,err)
return"",false
}
defer fid.Close()
fid.Seek(0,0)
b:=bufio.NewReader(fid)
for s,err:=b.ReadString('\n');err==nil||err==io.EOF;s,err= b.ReadString('\n'){
glog.V(debug).Infof("looking for a cid in '%s'\n",s)
if strings.HasPrefix(s,"Content-ID: <"){
s= s[len("Content-ID: <"):len(s)-2]
glog.V(debug).Infof("found a cid '%s'\n",s)
return s,true
}
if err==io.EOF{
break
}
}
return"",false
}





/*:174*/



/*178:*/


//line amail.w:2460

func readString(pfid*client.Fid,name string)(str string,err error){
glog.V(debug).Infof("readString: trying to open '%s'\n",name)
f,err:=pfid.Walk(name)
if err==nil{
err= f.Open(plan9.OREAD)
}
if err!=nil{
return
}
defer f.Close()
str,err= bufio.NewReader(f).ReadString('\n')
if err!=nil&&err!=io.EOF{
return
}
return str,nil
}



/*:178*/



/*179:*/


//line amail.w:2479

func readAll(pfid*client.Fid,name string,buf[]byte)([]byte,error){
glog.V(debug).Infof("readAll: trying to open '%s'\n",name)
f,err:=pfid.Walk(name)
if err==nil{
err= f.Open(plan9.OREAD)
}
if err!=nil{
return buf,err
}
defer f.Close()
b:=bufio.NewReader(f)
for s,err:=b.ReadString('\n');err==nil;s,err= b.ReadString('\n'){
if strings.HasSuffix(s,"\r\n"){
s= strings.TrimRight(s,"\r\n")
s+= "\n"
}
buf= append(buf,s...)
}
return buf,nil
}




/*:179*/



/*181:*/


//line amail.w:2539

func saveFile(src,dst string)error{
var err error
c:=exec.Command("9p","read",src)
f,err:=os.OpenFile(dst,os.O_WRONLY|os.O_CREATE|os.O_TRUNC,0600)
if err!=nil{
glog.Errorf("can't create a file '%s': %v\n",dst,err)
return err
}
defer f.Close()
c.Stdout= f
if err= c.Run();err!=nil{
glog.Errorf("can't run '9p': %v\n",err)
}
return err
}



/*:181*/



/*182:*/


//line amail.w:2557

func(this*message)fixFile(dir string)error{
src:=fmt.Sprintf("%d/%s",this.id,this.html)
dst:=fmt.Sprintf("%s/%d.html",dir,this.id)
df,err:=os.OpenFile(dst,os.O_WRONLY|os.O_CREATE|os.O_TRUNC,0600)
if err!=nil{
glog.Errorf("can't create a file '%s': %v\n",dst,err)
return err
}
defer df.Close()
fid,err:=this.box.fid.Walk(src)
if err==nil{
err= fid.Open(plan9.OREAD)
}
if err!=nil{
glog.Errorf("can't open to '%s': %v\n",src,err)
return err
}
defer fid.Close()
b:=bufio.NewReader(fid)
for s,err:=b.ReadString('\n');err==nil||err==io.EOF;s,err= b.ReadString('\n'){
p:=0
for b:=strings.Index(s[p:],"\"cid:");b!=-1;b= strings.Index(s[p:],"\"cid:"){
b+= p
e:=strings.Index(s[b+1:],"\"")
if e==-1{
break
}
e++
cid:=s[p+b+5:b+e]
glog.V(debug).Infof("cid: %s\n",cid)
if f,ok:=this.cids[cid];ok{
glog.V(debug).Infof("found a cid: %s, replace '%s' by '%s'\n",cid,s[b+1:b+e],f.name)
s= strings.Replace(s,s[b+1:b+e],f.name,1)
}else{
p= e
}
}
df.Write([]byte(s))
if err==io.EOF{
break
}
}
return err
}



/*:182*/


