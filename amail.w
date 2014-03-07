\input header

@** Introduction.

\.{Amail} is a mail client for \.{Acme} - an editor/window manager/shell.
It is supposed to be a replacement for \.{Mail} - the classic mail client for \.{Acme}.

For years I was being a user of \.{Opera} - a web browser with a mail client.
But a quality of the web browser of \.{Opera} was becoming low from a version to a version,
so I decided to change the web brower to a \.{Chromium}, but I didn't find a mail client for my requirements.

Few years ago I saw \.{Acme} and found it is very simple, but powerful and extremely extensible.
Yes, it is not perfect (nothing is perfect), but it is good enough, and I use it like a programming environment
(instead of \.{Emacs}).
I had known about \.{Mail} - a mail client for \.{Acme}, and a time to try it has come.

I have found \.{Mail} has some disadvantages (at least for me):
\yskip\item{$\bullet$}it doesn't have a support of threads
\yskip\item{$\bullet$}it doesn't have a support for read/unread messages
\yskip\item{$\bullet$}it doesn't have a navigation though mailboxes
\yskip\item{$\bullet$}it has a quite big loading time with big mailboxes.

I also prefer to view some messages in \.{html}-form (if any) with a possibility to open them in a web browser.

\.{Amail} is supposed to use with a conjunction with a \.{upas} - a mail filesystem supports \.{IMAP4} mail protocol.
I'm going to save a compatibility with \.{Mail} by commands.

For the moment \.{Amail} is working with \.{Acme} from \.{Plan 9 from User Space} (http://swtch.com/plan9port/).
I have some doubts \.{Amail} will work in \.{Plan9} without changes.

@** Implementation.
@c
@i license
import (
	@<Imports@>
)@#

type (
	@<Types@>
)@#

@<Constants@>@#

var (
	@<Variables@>
	debug glog.Level = 1
)@#

func main() {
	glog.V(debug).Infoln("main")
	defer glog.V(debug).Infoln("main is done")
	@<Parse command line arguments@>
	@<Try to open \.{mailfs}@>
	@<Subscribe on notifications of \.{plumber}@>
	@<Init root of \.{mailfs}@>
	@<Start a collector of message identifiers@>
	if len(flag.Args())>0 {
		@<Start a main message loop@>
		for _, name:=range flag.Args() {
			@<Create |box|@>
			@<Start a message loop for |box|@>
			@<Inform |box| to create a window@>
		}
	} else {
		@<Create the main window@>
		@<Start a main message loop@>
		go func(){
			@<Enumerating of mailboxes@>
		}()
	}
	@<Process events are specific for |boxes|@>
}

@* Exiting.

At first we should quit correctly.
So a channel |exit| is defined. All goroutines should wait for a message from |exit|.
@<Variables@>=
exit chan bool=make(chan bool)

@
@<On exit?@>=
	case <-exit:
	glog.V(debug).Infoln("on exit!")
	return

@
@<Exit!@>=
glog.V(debug).Infoln("exit!")
close(exit)

@ We have to quit when all window of mailboxes and main window are closed.
|wcount| contains a count of mailboxes's windows.
|wch| is a channel to manipulate of |wcount|.
When the main window is closed, the program has to exit immediately.
@<Variables@>=
wch chan int=make(chan int, 100)
wcount int

@ When |wcount==0|, the program quits.
@<Processing of other common channels@>=
	case i:=<-wch:
		wcount+=i
		if wcount==0 {
			@<Exit!@>
			return
		}

@
@<Increase the windows count@>=
glog.V(debug).Infoln("increase the windows count")
wch<-1

@
@<Decrease the windows count@>=
glog.V(debug).Infoln("decrease the windows count")
wch<-(-1)


@* Parsing command line arguments.
@<Variables@>=
shownew		bool
showthreads	bool
levelmark	string
newmark		string
skipboxes	[]string

@
@<Imports@>=
"flag"
"fmt"
"os"
"strings"
"sort"

@
@<Parse command line arguments@>=
{
	glog.V(debug).Infoln("parsing command line arguments")
	var skip string
	flag.BoolVar(&shownew, "new", true, "show new messages only")
	flag.BoolVar(&showthreads, "threads", true, "show threads of messages")
	flag.StringVar(&skip, "skip", "", "boxes to be skiped, separated by comma")
	flag.StringVar(&levelmark, "levelmark", "+", "mark of level for threads")
	flag.StringVar(&newmark, "newmark", "(*)", "mark of new messages")
	flag.Usage=func() {
		fmt.Fprintf(os.Stderr, "Mail client for Acme programming environment\n")
		fmt.Fprintf(os.Stderr, "Usage: %s [options] [<mailbox 1>]...[<mailbox N>]\n", os.Args[0])
		fmt.Fprintln(os.Stderr, "Options:")
		flag.PrintDefaults()
	}
	flag.Parse()
	@<Check |levelmark| and |newmark|@>
	if len(skip)>0 {
		skipboxes=strings.Split(skip, ", ")
		sort.Strings(skipboxes)
		glog.V(debug).Infof("these mailboxes will be skipped: %v\n", skipboxes)

	}
}

@
@<Imports@>=
"unicode"
"unicode/utf8"

@ |levelmark| shouldn't have ending digit and |newmark| shouldn't have leading digit,
because the digits change a message id.
@<Check |levelmark| and |newmark|@>=
glog.V(debug).Infoln("checking of levelmark and newmark")
if r, _:=utf8.DecodeLastRuneInString(levelmark); unicode.IsDigit(r) {
	fmt.Fprintln(os.Stderr, "last symbol of level mark shouldn't be a digit")
	os.Exit(1)
}
if r, _:=utf8.DecodeRuneInString(newmark); unicode.IsDigit(r) {
	fmt.Fprintln(os.Stderr, "first symbol of new mark shouldn't be a digit")
	os.Exit(1)
}

@* Mounting of the \.{Acme} filesystem.
@<Imports@>=
"bitbucket.org/santucco/goplan9-clone/plan9/client"
"github.com/golang/glog"

@
@<Variables@>=
fsys	*client.Fsys
rfid	*client.Fid
srv		string="mail"

@
@<Try to open \.{mailfs}@>=
{
	glog.V(debug).Infoln("try to open mailfs")
	var err error
	if fsys, err=client.MountService(srv); err!=nil {
		glog.Errorf("can't mount mailfs: %v\n", err)
		os.Exit(1)
	}
}


@* Enumeration of mailboxes.
@ Let's make a structure of a mailbox.
@<Types@>=
mailbox struct {
	name	string
	@<Rest of |mailbox| members@>
}

mailboxes []*mailbox

message struct {
	id	int
	@<Rest of |message| members@>
}

messages []*message


@ |boxes| contains all message boxes are enumerated and sorted
@<Variables@>=
boxes	mailboxes

@ The |mailbox| structure has to be extended a bit:
\yskip\item{$\bullet$}|all| is a list of all messages in the box;
\yskip\item{$\bullet$}|unread| is a list of unread messages in the box;
\yskip\item{$\bullet$}|mch| is a channel to manipulate of |all| and |unread|;
\yskip\item{$\bullet$}|dch| is a channel to inform the box a message has been deleted.

@<Rest of |mailbox| members@>=
all		messages
unread	messages
mch		chan int
dch 	chan int

@
@<Rest of initialization of |mailbox|@>=
mch:make(chan int, 100), @/
dch:make(chan int, 100), @/


@ Four global channels for synchronious mails counting should be defined:
\yskip\item{$\bullet$}|mch| is a channel receives info about a message from \.{plumber};
\yskip\item{$\bullet$}|dch| is a channel receives info about deleted message from \.{plumber};
\yskip\item{$\bullet$}|bch| is a channel receives info about new boxes;
\yskip\item{$\bullet$}|rfch| is a channel receives info about a box should be refreshed in the main window.

@<Variables@>=
mch=make(chan *struct{name string; id int}, 100)
dch=make(chan *struct{name string; id int}, 100)
bch=make(chan string, 10)
rfch=make(chan *mailbox, 100)

@ A slice of enumerated mailboxes should be sorted.
A few methods have to be implemented for |mailboxes| to have an ability to sort of them
@c
func (this mailboxes) Len() int {
	return len(this)
} @#

func (this mailboxes) Less(i, j int) bool {
	return this[i].name < this[j].name
} @#

func (this mailboxes) Swap(i, j int) {
	t:=this[i]
	this[i]=this[j]
	this[j]=t
} @#

@ Here we open the root of mailfs
@<Init root of \.{mailfs}@>=
	glog.V(debug).Infoln("initialization of root of mailfs")
	var err error
	rfid, err=fsys.Walk(".")
	if err!=nil {
		glog.Errorf("can't open mailfs: %v\n", err)
		os.Exit(1)
	}
	defer rfid.Close()


@ Here we read all directory names.
@<Enumerating of mailboxes@>=
{
	glog.V(debug).Infoln("enumerating of mailboxes")
	fi, err:=rfid.Dirreadall()
	if err!=nil {
		glog.Errorf("can't read mailfs: %v\n", err)
		@<Exit!@>
		return
	}
	for _, f:=range fi {
		if f.Mode&plan9.DMDIR==plan9.DMDIR {
			name:=f.Name
			@<Add a mailbox...@>
		}
	}
	glog.V(debug).Infoln("enumerating of mailboxes is done")
}

@ Names of directories are sent in |bch|
@<Add a mailbox with |name|@>=
glog.V(debug).Infof("send a mailbox '%s' to put in the list\n", name)
bch<-name

@ |newMessage| is a method of |mailbox| to fill a message with |id|.
@<Imports@>=
"io"
"bufio"

@
@<Rest of |message| members@>=
unread bool
box	*mailbox

@ |newMessage| creates |msg| and fills its fields from |"info"| file.
|"flags"| are parsed to detect the message is new.
@c
func (this *mailbox) newMessage(id int) (msg *message, unread bool, err error) {
	glog.V(debug).Infof("newMessage: trying to open '%d/info'\n", id)
	f, err:=this.fid.Walk(fmt.Sprintf("%d/info", id))
	if err==nil {
		err=f.Open(plan9.OREAD)
	}
	if err!=nil {
		glog.Errorf("can't open to '%s/%d/info': %s\n", this.name, id, err)
		return
	}
	defer f.Close()
	msg=&message{id:id, box:this, @<Rest of initialization of |message|@>}
	b:=bufio.NewReader(f)
	unread=true
	glog.V(debug).Infof("newMessage: reading and parsing of a content of '%d/info'\n", id)
	for s, err:=b.ReadString('\n'); err==nil;  s, err=b.ReadString('\n') {
		if s[len(s)-1]=='\n' {
			s=s[:len(s)-1]
		}
		if strings.HasPrefix(s, "flags ") {
			if strings.Index(s, "seen")>=0 {
				unread=false
			}
			continue
		}
		@<Read other fields of a message@>
	}
	msg.unread=unread
	return

}

@* Subscription on notifications about new messages.

@<Imports@>=
"bitbucket.org/santucco/goplumb"
"bitbucket.org/santucco/goplan9-clone/plan9"

@ Here a subscription on |"seemail"| port of \.{plumber} is made. The messages is checked for |filetype=="mail"| and
|"mailtype"| are existing. In case a new mail message we send a name of a mailbox an an id of the message in |mch|,
in case of a mail message is deleted - in |dch|.
@<Subscribe on notifications of \.{plumber}@>=
{
	glog.V(debug).Infoln("trying to open 'seemail' plumbing port")
	if sm, err:=goplumb.Open("seemail", plan9.OREAD); err!=nil {
		glog.Errorf("can't open plumb/seemail: %s\n", err)
	} else {
		sch, err:=sm.MessageChannel(0)
		if err!=nil {
			glog.Errorf("can't get message channal for plumb/seemail: %s\n", err)
		} else {
			go func() {
				defer sm.Close()
				defer glog.V(debug).Infoln("plumbing goroutine is done")
				for {
					select {
						@<On exit?@>
						case m, ok:=<-sch:
							if !ok {
								glog.Warningln("it seems plumber has finished")
								sch=nil
								return
							}
							glog.V(debug).Infof("a plumbing message has been received: %v\n", m)
							if m.Attr["filetype"]!="mail" {
								glog.Warningln("attribute 'filetype' is not 'mail'")
								continue
							}
							v, ok:=m.Attr["mailtype"]
							if !ok {
								glog.Warningln("can't find 'mailtype' attribute")
								continue
							}
							s:=strings.TrimLeft(string(m.Data), "Mail/")
							n:=strings.LastIndex(s, "/")
							if n==-1 {
								glog.Warning("can't found a number of message in '%s'\n", s)
								continue
							}
							num, err:=strconv.Atoi(s[n+1:])
							if err!=nil {
								glog.Error(err)
								continue
							}
							if v=="new" {
								glog.V(debug).Infof("'%d' is a new message in the '%s' mailbox\n", num, s[:n])
								mch<-&struct{name string; id int}{name:s[:n], id:num}
							} else if v=="delete" {
								glog.V(debug).Infof("'%d' is a deleted message in the '%s' mailbox\n", num, s[:n])
								dch<-&struct{name string; id int}{name:s[:n], id:num}
							}
					}
				}
			} ()
		}
	}
}

@ Here a subscription on |"sendmail"| port of \.{plumber} is made.
@<Subscribe on notifications of \.{plumber}@>=
{
	glog.V(debug).Infoln("trying to open 'sendmail' plumbing port")
	if sm, err:=goplumb.Open("sendmail", plan9.OREAD); err!=nil {
		glog.V(debug).Infof("can't open plumb/sendmail: %s\n", err)
	} else {
		sch, err:=sm.MessageChannel(0)
		if err!=nil {
			glog.Errorf("can't get message channal for plumb/sendmail: %s\n", err)
		} else {
			go func() {
				defer sm.Close()
				defer glog.V(debug).Infoln("plumbing goroutine is done")
				for {
					select {
						@<On exit?@>
						case m, ok:=<-sch:
							if !ok {
								glog.Warningln("it seems plumber has finished")
								sch=nil
								return
							}
							glog.V(debug).Infof("a plumbing message has been received: %v\n", m)
							var msg *message
							@<Create a new message window@>
							name:=fmt.Sprintf("Amail/New")
							@<Print the |name| for window |w|@>
							addr:=fmt.Sprintf("To: %s\n\n", string(m.Data))
							w.Write([]byte(addr))
							@<Append common signature@>
					}
				}
			} ()
		}
	}
}

@* The main message loop.

Via |bch| names of new mailboxes are received, the mailboxes is created and processed.
Via |mch| and |dch| messages about new and deleted messages are received, the corresponding mailboxes are found and
the messages identifiers are send in the corresponding channels of the mailboxes.
@<Process events are specific for |boxes|@>=
glog.V(debug).Infoln("process events are specific for the list of mailboxes")
for {
	select {
		@<On exit?@>
		case name:=<-bch:
			@<Continue if the box |name| should be skiped@>
			@<Create |box|@>
			@<Send a signal to refresh all mailboxes@>
			@<Start a message loop for |box|@>
		case d:=<-mch:
			name:=d.name
			@<Looking for a |name| mailbox...@>
			glog.V(debug).Infof("sending '%d' to add in the '%s' mailbox\n", d.id, boxes[i].name)
			boxes[i].mch<-d.id
		case d:=<-dch:
			name:=d.name
			@<Looking for a |name| mailbox...@>
			glog.V(debug).Infof("sending '%d' to delete from the '%s' mailbox\n", d.id, boxes[i].name)
			boxes[i].dch<-d.id
		@<Processing of other common channels@>
	}
}


@ This is a message loop for main window. It reads and processes messages from different channels.

A pointer to a mailbox |b| is received from |rfch|. In case |b==nil| we should print a state of all mailboxes
or state of |b| otherwise.
@<Start a main message loop@>=
go func() {
	glog.V(debug).Infoln("start a main message loop")
	defer glog.V(debug).Infoln("main message loop is done")
	for {
		select {
			@<On exit?@>
			case b:=<-rfch:
				if b==nil {
					@<Print all mailboxes@>
				} else {
					@<Refresh main window...@>
				}
			@<Processing of other channels@>
		}
	}
}()

@
@<Create |box|@>=
glog.V(debug).Infof("creating a '%s' mailbox\n", name)
box:=&mailbox{name:name, @<Rest of initialization of |mailbox|@>}
boxes=append(boxes, box)
sort.Sort(boxes)

@ |mailboxes.Search| finds a mailbox with |name| and returns a position of the mailbox in the list and |true| or
a position where the box can be inserted and |false|
@c
func (this mailboxes) Search(name string) (int, bool) {
	pos:=sort.Search(len(this), @t\1@>@/
		 func(i int) bool {return this[i].name>=name} @t\2@>);
	if pos!=len(this) && this[pos].name==name {
		return pos, true
	}
	return pos, false
}

@
@<Looking for a |name| mailbox, storing an index of the mail box was found in |i|, continue if not found@>=
glog.V(debug).Infof("looking for a '%s' mailbox\n", name)
i, ok:=boxes.Search(name)
if !ok {
	@<Continue if the box |name| should be skiped@>
	glog.Warningf("can't find message box '%s'\n", name)
	continue
}


@* The main window.
@<Imports@>=
"bitbucket.org/santucco/goacme"

@
@<Variables@>=
mw		*goacme.Window
ech		<-chan *goacme.Event

@
@<Create the main window@>=
glog.V(debug).Infoln("creating the main window")
defer goacme.DeleteAll()

var err error
if mw, err=goacme.New(); err!=nil {
	glog.Errorf("can't create a window: %v\n", err)
	os.Exit(1)
}
name:="Amail"
w:=mw
@<Print the |name| for window |w|@>
if ech, err=mw.EventChannel(0, goacme.Look|goacme.Execute); err!=nil {
	glog.Errorf("can't open an event channel of the window %v\n", err)
	os.Exit(1)
}
@<Write a tag of main window@>
@<Increase the windows count@>

@
@<Constants@>=
const mailboxfmt="%-30s\t%10d\t%10d\n"
const mailboxfmtprc="%-30s\t%10d\t%10d\t%d%%\n"
const wholefile="0,$"

@ Here we clean up the main window and print states of all mailboxes.
@<Print all mailboxes@>=
if mw!=nil {
	glog.V(debug).Infoln("printing of the mailboxes")
	if err:=mw.WriteAddr(wholefile); err!=nil {
		glog.Errorf("can't write '%s' to 'addr' file: %s\n", wholefile, err)
	} else if data, err:=mw.File("data"); err!=nil {
		glog.Errorf("can't open 'data' file: %s\n", err)
	} else {
		for _, v:=range boxes {
			if v.total==len(v.all) {
				data.Write([]byte(fmt.Sprintf(mailboxfmt, v.name, len(v.unread), len(v.all))))
			} else if v.total!=0 && len(v.all)*100/v.total>0 {
				data.Write([]byte(fmt.Sprintf(mailboxfmtprc, v.name, len(v.unread), len(v.all), len(v.all)*100/v.total)))
			} else {
				data.Write([]byte(fmt.Sprintf(mailboxfmt, v.name, 0, 0)))
			}
		}
	}
	w:=mw
	@<Set window |w| to clean state@>
	@<Go to top of window |w|@>
}


@ Let's add processing of evens from the main window.
Here events from the main window are processed.
@<Processing of other channels@>=
case ev, ok:=<-ech:
	glog.V(debug).Infof("an event from main window has been received: %v\n", ev)
	if !ok {
		ech=nil
		@<Decrease the windows count@>
		return
	}
	if (ev.Type&goacme.Execute)==goacme.Execute {
		switch ev.Text {
			case "ShowNew":
				shownew=true
			case "ShowAll":
				shownew=false
			case "ShowThreads":
				showthreads=true
			case "ShowPlain":
				showthreads=false
			case "Del":
				mw.UnreadEvent(ev)
				mw.Close()
				mw=nil
				@<Exit!@>
				return
			case "debug":
				debug=0
				continue
			case "nodebug":
				debug=1
				continue
			default:
				mw.UnreadEvent(ev)
				continue
		}
		@<Write a tag of main window@>
		continue
	} else if (ev.Type&goacme.Look)==goacme.Look {
		name:=ev.Text
		// a box name can contain spaces
		if len(ev.Arg)>0 {
			name+=" "+ev.Arg
		}
		name=strings.TrimSpace(name)
		if i, ok:=boxes.Search(name); ok {
			box:=boxes[i]
			@<Inform |box| to create a window@>
			continue
		}
	}
	mw.UnreadEvent(ev)


@ If not all messages are counted, the refresh of state of mailbox in the main window will be processed every percent
of counted messages.

@<Refresh main window for a box |b|@>=
glog.V(debug).Infof("refreshing main window for the '%s' mailbox, len(all): %d, total: %d\n", b.name, len(b.all), b.total)
if mw!=nil {
	if len(b.all)!=b.total && b.total/100!=0 && len(b.all)%(b.total/100)!=0 {
		continue
	}

	if err:=mw.WriteAddr("0/^%s.*\\n/", escape(b.name)); err!=nil {
		glog.V(debug).Infof("can't write to 'addr' file: %s\n", err)
		continue
	}

	if data, err:=mw.File("data"); err !=nil {
		glog.Errorf("can't open 'data' file: %s\n", err)
	} else if len(b.all)==b.total {
		if _, err:=data.Write([]byte(fmt.Sprintf(mailboxfmt, b.name, len(b.unread), len(b.all))));
			err!=nil {
			glog.Errorf("can't write to 'data' file: %s\n", err)
			continue
		}
		w:=mw
		@<Set window |w| to clean state@>
		@<Go to top of window |w|@>
	} else if _, err:=data.Write([]byte(fmt.Sprintf(mailboxfmtprc,
													b.name,
													len(b.unread),
													len(b.all),
													len(b.all)*100/b.total)));err!=nil {
		glog.Errorf("can't write to 'data' file: %s\n", err)
	}
}

@
@<Imports@>=
"strconv"

@
@<Rest of |mailbox| members@>=
fid		*client.Fid
total	int

@ Here messages of a mailbox are counted. If some directories are not numbers, they are supposed to be mailboxes
and its names are sent to |bch| to start of counting of the new mailbox.
New messages are counted here too.
The enumeration of the messages is started from the end of the list, because new messages have bigger numbers.
To avoid of unresponding main window the counting is made in |default| branch of |select|.

@<Count of messages in a box@>=
{
	glog.V(debug).Infof("counting of messages in the '%s' mailbox\n", box.name)
	var err error
	box.fid, err=rfid.Walk(box.name)
	if err!=nil {
		glog.Errorf("can't walk to '%s': %v", box.name, err)
		return
	}
	defer box.fid.Close()
	fs, err:=box.fid.Dirreadall()
	if err!=nil {
		glog.Errorf("can't read a mailbox '%s': %s", box.name, err)
		return
	}
	box.total=len(fs)-2
	box.all=make(messages, 0, box.total)
	for i:=len(fs)-1; i>=0; {
		select {
			@<On exit?@>
			@<Processing of other |box| channels@>
			default:
				d:=fs[i]
				i--
				if (d.Mode&plan9.DMDIR)!=plan9.DMDIR {
					continue
				}
				id, err:=strconv.Atoi(d.Name)
				if err!=nil { // it seems this is a mailbox
					// decrease a total number of messages
					box.total--
					name:=box.name+"/"+d.Name
					@<Add a mailbox...@>
					continue
				}
				if msg, new, err:=box.newMessage(id); err==nil {
					if new {
						@<Add |msg| to |unread|@>
					}
					@<Add |msg| to |all|@>
				} else {
					glog.V(debug).Infof("can't create a new '%d' message in the '%s' mailbox: %v\n", id, box.name, err)
					box.total--
					continue
				}
				@<Send |box| to refresh the main window@>
				@<Inform |box| to print counting messages@>
		}
	}
	@<Inform |box| to print the rest of counting messages@>
}

@
@<Set window |w| to clean state@>=
if w!=nil{
	glog.V(debug).Infoln("setting the window to clean state")
	if err:=w.WriteCtl("clean"); err!=nil {
		glog.Errorf("can't write to 'ctl' file: %s\n", err)
	}
}

@
@<Set window |w| to dirty state@>=
if w!=nil{
	glog.V(debug).Infoln("setting the window to dirty state")
	if err:=w.WriteCtl("dirty"); err!=nil {
		glog.Errorf("can't write to 'ctl' file: %s\n", err)
	}
}

@ |escape| escapes the regex specific charactets
@c
func escape(s string) (res string) {
	for _, v:=range s {
		if strings.ContainsRune("\\/[].+?()*^$", v) {
			res+="\\"
		}
		res+=string(v)
	}
	return res
}

@ If |name| contains spaces, they will be replaced by underlines.
@<Print the |name| for window |w|@>=
glog.V(debug).Infoln("printing a name for a window")
if err:=w.WriteCtl("name %s", strings.Replace(name, " ", "␣", -1)); err!=nil {
	glog.Errorf("can't write to 'ctl' file: %s\n", err)
}

@
@<Continue if the box |name| should be skiped@>=
glog.V(debug).Infoln("continue if the box should be skiped")
if i:=sort.SearchStrings(skipboxes, name); i!=len(skipboxes) && skipboxes[i]==name {
	continue
}

@ |messages.Search| finds a message with |id| and returns a position of the message in the list and |true| or
a position where the message can be inserted and |false|
@c
func (this messages) Search(id int) (int, bool) {
	pos:=sort.Search(len(this), func(i int) bool {return this[i].id<=id});
	if pos!=len(this) && this[pos].id==id {
		return pos, true
	}
	return pos, false
}

@ |messages.Insert| inserts a message |msg| in position |pos|
@c
func (this *messages) Insert(msg *message, pos int) {
	*this=append(*this, nil)
	copy((*this)[pos+1:], (*this)[pos:])
	(*this)[pos]=msg
}

@ |messages.SearchInsert| inserts a message |msg| and returns
a position of the message in the list and |true| or
a position where the message already exists and |false|
@c
func (this *messages) SearchInsert(msg *message) (int, bool){
	pos, ok:=this.Search(msg.id)
	if ok {
		return pos, false
	}
	this.Insert(msg, pos)
	return pos, true
}

@ |messages.Delete| deletes a message at |pos| position and returns
a pointer to the message is removed and |true| if the message is deleted,
|false| otherwise
@c
func (this *messages) Delete(pos int) (*message, bool) {
	if pos<0 || pos>len(*this)-1 {
		return nil, false
	}
	msg:=(*this)[pos]
	*this=append((*this)[:pos], (*this)[pos+1:]...)
	return msg, true
}

@ |messages.DeleteById| deletes a message with |id| and returns
a pointer to the message is removed and |true| if the message is deleted,
|false| otherwise
@c
func (this *messages) DeleteById(id int) (*message, bool) {
	pos, ok:=this.Search(id)
	if !ok {
		return nil, false
	}
	return this.Delete(pos)
}

@* The message loop for a mailbox.
@<Start a message loop for |box|@>=
go box.loop()

@
@c
func(box *mailbox) loop() {
	glog.V(debug).Infof("start a message loop for the '%s' mailbox\n", box.name)
	counted:=false
	pcount:=0
	ontop:=false
	@<Count of messages in a box@>
	counted=true
	@<Write a tag of |box| window@>
	if box.threadMode() {
		@<Inform |box| to print messages@>
	}
	defer glog.V(debug).Infof("a message loop of the '%s' mailbox is done\n", box.name)
	for {
		select {
			@<On exit?@>
			@<Processing of other |box| channels@>
		}
	}
}

@ Here new messages of |box| are processed.
If |box| shows a particular thread the message should be printed only if it is in the thread.
Also the message should be printed in other boxes with the thread.
A new message can be in a thread so we have to send
@<Processing of other |box| channels@>=
case id:=<-box.mch:
	glog.V(debug).Infof("'%d' should be added to the '%s' mailbox\n", id, box.name)
	msg, new, err:=box.newMessage(id)
	if err!=nil {
		continue
	}
	if new {
		@<Add |msg| to |unread|@>
	}
	box.total++
	@<Add |msg| to |all|@>
	@<Print |msg| at exact positon@>
	if !box.thread {
		if box.threadMode() {
			@<Get |root| of |msg|@>
			var msgs messages
			src:=append(messages{}, root)
			@<Make a full thread in |msgs| for |root|@>
			@<Inform |box| to print |msgs|@>
		} else {
			@<Inform |box| to print |msg|@>
		}
	}
	@<Send |box| to refresh the main window@>


@ Here deleted messages of |box| are processed.
@<Processing of other |box| channels@>=
case id:=<-box.dch:
	glog.V(debug).Infof("'%d' should be deleted from the '%s' mailbox\n", id, box.name)
	@<Delete a message with |id|@>

@ |deleted| points out the message is marked to delete.
@<Rest of |message| members@>=
deleted bool

@
@<Add |msg| to |unread|@>=
{
	glog.V(debug).Infof("addition of the '%d' message to the list of unread messages of the '%s' mailbox\n",
		msg.id, box.name )
	box.unread.SearchInsert(msg)
}

@
@<Add |msg| to |all|@>=
{
	glog.V(debug).Infof("addition of the '%d' message to the list of all messages of the '%s' mailbox\n",
		msg.id, box.name )
	box.all.SearchInsert(msg)
}

@
@<Delete a message with |id|@>=
if i, ok:=box.all.Search(id); ok {
	msgs:=append(messages{}, box.all[i])
	@<Delete a message at position |i|@>
	@<Send deleted |msgs|@>
}

@ Here we delete a message from |box.all| and |box.unread|, decrease |total| and |deleted| counts,
send the message id to clean a thread links and send a signal to refresh the main window.
@<Delete a message at position |i|@>=
{
	if msg, ok:=box.all.Delete(i); ok {
		glog.V(debug).Infof("deleting the '%d' message from the '%s' mailbox\n", msg.id, box.name)
		box.unread.DeleteById(msg.id)
		box.total--
		if msg.deleted {
			box.deleted--
		}
		@<Clean up |msg|@>
		@<Send |box| to refresh the main window@>
	}
}

@
@c
func (box *mailbox) threadMode() bool {
	return box.thread || box.showthreads && !box.shownew
}

@ Here we make a snapshot of |box| state and send it to |rfch|
@<Send |box| to refresh the main window@>=
glog.V(debug).Infof("sending a snapshot of the '%s' mailbox to refresh the main window\n", box.name)
b:=*box
rfch<-&b

@
@<Send a signal to refresh all mailboxes@>=
glog.V(debug).Infoln("sending a signal to refresh all mailboxes")
rfch<-nil

@ Let's add some members to |mailbox|. |shownew| and |showthreads| is a copy of global corresponding flags.

|ech| is a channel of events from |box|'s window.

|w| is a |box|'s window.

|cch| is a channel receives the signal to create |box|'s window.

@<Rest of |mailbox| members@>=
shownew		bool
showthreads	bool
ech 		<-chan *goacme.Event
w			*goacme.Window
cch			chan bool

@
@<Rest of initialization of |mailbox|@>=
shownew:shownew, @/
showthreads:showthreads, @/
cch:make(chan bool, 100),

@
@<Inform |box| to create a window@>=
glog.V(debug).Infof("inform the '%s' mailbox to create a window\n", box.name)
box.cch<-true

@ Here we are waiting for a signal to create |box|'s window and create it.
@<Processing of other |box| channels@>=
case <-box.cch:
	glog.V(debug).Infof("a signal to create the '%s' mailbox window has been received\n", box.name)
	if box.w==nil {
		box.shownew=shownew
		box.showthreads=showthreads
		box.thread=false
		@<Create a window for the box@>
		@<Inform |box| to print messages@>
		@<Increase the windows count@>
	} else {
		glog.V(debug).Infof("a window of the '%s' mailbox already exists, just show it\n", box.name)
		box.w.WriteCtl("dot=addr\nshow")
	}

@
@<Create a window for the box@>=
glog.V(debug).Infof("creation a window for the '%s' mailbox\n", box.name)
var err error
if box.w, err=goacme.New(); err!=nil {
	glog.Errorf("can't create a window: %v\n", err)
	os.Exit(1)
}
if box.ech, err=box.w.EventChannel(0, goacme.Look|goacme.Execute); err!=nil {
	glog.Errorf("can't open an event channel of the window %v\n", err)
	os.Exit(1)
}
@<Write a name of |box| window@>
@<Write a tag of |box| window@>


@ |thread| flag points out the |box|'s window shows a particular thread of messages.
@<Rest of |mailbox| members@>=
thread bool

@ Processing of events from the box's window
@<Processing of other |box| channels@>=
case ev, ok:=<-box.ech:
	glog.V(debug).Infof("an event has been received from the '%s' mailbox window: %v\n", box.name, ev)
	if !ok {
		box.ech=nil
		continue
	}
	if (ev.Type&goacme.Execute)==goacme.Execute {
		switch ev.Text {
			case "Del":
				@<Clean window-specific stuff@>
				box.w.Del(true)
				box.w.Close()
				box.w=nil
				@<Decrease the windows count@>
				continue
			case "ShowNew":
				box.thread=false
				box.shownew=true
			case "ShowAll":
				if box.showthreads && !counted {
					continue
				}
				box.thread=false
				box.shownew=false
			case "ShowThreads":
				if !counted {
					continue
				}
				box.showthreads=true
				if box.shownew==true {
					@<Write a tag of |box| window@>
					continue
				}
			case "ShowPlain":
				box.showthreads=false
				if box.shownew==true {
					@<Write a tag of |box| window@>
					continue
				}
			case "Thread":
				if !counted {
					continue
				}
				var msg *message
				if len(ev.Arg)==0 {
					@<Get a pointer |msg| to current message@>
				} else if num, err:=strconv.Atoi(strings.TrimSpace(ev.Arg)); err!=nil {
					continue
				} else if p, ok:=box.all.Search(num); ok {
					msg=box.all[p]
				}
				if msg!=nil {
					box.thread=true
					@<Write a tag of |box| window@>
					@<Clean |box| window@>
					@<Clean window-specific stuff@>
					@<Inform |box| to print a full thread with |msg|@>
				}
				continue
			case "Delmesg":
				@<Mark to delete messages@>
				continue
			case "UnDelmesg":
				@<Unmark messages@>
				continue
			case "Put":
				@<Delete messages@>
				continue
			case "Mail":
				var msg *message
				@<Create a new message window@>
				name:=fmt.Sprintf("Amail/%s/New", box.name)
				@<Print the |name| for window |w|@>
				@<Append |box|-specific signature@>
				continue
			case "Seen":
				@<Mark messages as seen@>
				continue
			case "Search":
				glog.V(debug).Infof("search argument: '%s'\n", ev.Arg)
				@<Search messages@>
				continue
			default:
				box.w.UnreadEvent(ev)
				continue
		}
		@<Write a name of |box| window@>
		@<Write a tag of |box| window@>
		@<Clean |box| window@>
		@<Set window |w| to clean state@>
		@<Clean window-specific stuff@>
		@<Inform |box| to print messages@>
		continue
	} else if (ev.Type&goacme.Look)==goacme.Look {
		@<Create |msgs|@>
		if (ev.Type&goacme.Tag)==goacme.Tag {
			s:=ev.Text
			@<Read a message number@>
		} else {
			@<Open selected messages@>
		}
		if len(msgs)!=0 {
			@<Send |msgs| for viewing@>
			continue
		}
	}
	box.w.UnreadEvent(ev)

@ Several messages can be selected to open. The address in |ev| will be inspected instead of |ev.Text|,
because a size of the selected messages can be more that 256 symbols. The address will send to |"addr"| file
of the |box|'s window and then symbols will be read from |"xdata"| file.

@<Open selected messages@>=
glog.V(debug).Infof("event: %v\n", ev)
if err:=box.w.WriteAddr("#%d,#%d", ev.Begin, ev.End); err!=nil {
	glog.Errorf("can't write to 'addr': %s\n", err)
} else
	@<Read message numbers@>

@
@<Read message numbers@>=
if xdata, err:=box.w.File("xdata"); err!=nil {
	glog.Errorf("can't open 'xdata' file: %s\n", err)
} else {
	b:=bufio.NewReader(xdata)
	for s, err:=b.ReadString('\n'); err==nil || err==io.EOF; s, err=b.ReadString('\n') {
		@<Read a message number@>
		if err==io.EOF {
			break
		}
	}
}


@ A message path can contain not only a number but a mailbox name too. So we have to parse an input string
to separate the name and the number. In any case the message will be opened via the main loop.
@<Read a message number@>=
{
	glog.V(debug).Infof("looking a message number in '%s'\n", s)
	s=strings.TrimLeft(s, levelmark+deleted)
	f:=strings.Split(s, "/")
	glog.V(debug).Infof("parts of message path: '%v'\n", f)
	id:=0
	for i, v:=range f {
		var err error
		if id, err=strconv.Atoi(strings.TrimRight(v, newmark)); err==nil {
			name:=box.name
			if i>0 {
				name=strings.Join(f[:i], "/")
				glog.V(debug).Infof("the message number is '%d' in the '%s' mailbox\n", id, name)
			}
			@<Add a |id| message to |msgs|@>
			break
		}
	}
}

@
@<Types@>=
msgmap map[string][]int

@ Let's do an universal approach to handle actions with groups of messages
instead of making a separate channel for every action.

@<Types@>=
action int

@
@<Constants@>=
const (
	view	action = iota
	@<Other kinds of actions@>
)

@ A channel to do actions with groups of messages.
@<Variables@>=
ach=make(chan *struct{m msgmap; a action}, 100)

@
@<Create |msgs|@>=
msgs:=make(msgmap)

@
@<Add a |id| message to |msgs|@>=
glog.V(debug).Infof("adding the '%d' of the '%s' mailbox\n", id, name)
msgs[name]=append(msgs[name], id)

@ Let's add a processing of |ach| to the main thread
@<Processing of other common channels@>=
	case d:=<-ach:
		if d.m==nil {
			continue
		}
		for name, ids:=range d.m {
			@<Looking for a |name| mailbox...@>
			boxes[i].ach<-&struct{ids []int; a action}{ids, d.a}
		}

@
@<Send |msgs| for viewing@>=
ach<-&struct{m msgmap; a action}{msgs, view}

@
@<Get a pointer |msg| to current message@>=
glog.V(debug).Infof("getting a pointer to current message in the window of the '%s' mailbox\n", box.name)
num:=0
if err:=box.w.WriteCtl("addr=dot"); err!=nil {
	glog.Errorf("can't write to 'ctl': %s\n", err)
} else if err:=box.w.WriteAddr("-/^/"); err!=nil {
	glog.V(debug).Infof("can't write to 'addr': %v\n", err)
} else if err:=box.w.WriteAddr("/[0-9]+(%s)?\\//", escape(newmark)); err!=nil {
	glog.V(debug).Infof("can't write to 'addr': %s\n", err)
} else  if data, err:=box.w.File("data"); err!=nil {
	glog.Errorf("can't open 'data' file: %s\n", err)
} else if str, err:=bufio.NewReader(data).ReadString('/'); err!=nil {
	glog.Errorf("can't read from 'data' file: %s\n", err)
} else if _, err:=fmt.Sscanf(strings.TrimLeft(str, levelmark), "%d", &num); err==nil {
	glog.V(debug).Infof("current message: %d(%s)\n", num, str)
	if p, ok:=box.all.Search(num); ok {
		msg=box.all[p]
	}
} else {
	glog.V(debug).Infof("can't get a current message from: %s\n", str)
}

@
@<Variables@>=
deleted="(deleted)-"

@
@<Compose a header of |msg|@>=
glog.V(debug).Infof("compose a header of the '%d' message of the '%s' mailbox\n", msg.id, box.name)
buf=append(buf, fmt.Sprintf("%s%s%d%s/\t%s\t%s\n\t%s\n", @t\1@>@/
	func() string{if msg.deleted {return deleted};return ""}(), @/
	func() string{if msg.box!=box {return fmt.Sprintf("%s/", msg.box.name)};return ""}(), @/
	msg.id, @/
	func() string{if msg.unread {return newmark};return ""}(), @/
	msg.from, @/
	msg.date, @/
	msg.subject )...@t\2@>)


@
@<Imports@>=
	"time"

@
@<Rest of |message| members@>=
from	string
date	time.Time
subject string

@
@<Read other fields of a message@>=
if strings.HasPrefix(s, "from ") {
	msg.from=s[len("from "):]
	msg.from=strings.Replace(msg.from, "'' ", "", -1)
	continue
}
var unixdate int64
if _, err:=fmt.Sscanf(s, "unixdate %d", &unixdate); err==nil {
	msg.date=time.Unix(unixdate, 0)
	continue
}
if strings.HasPrefix(s, "subject ") {
	msg.subject=s[len("subject "):]
	continue
}

@
@<Go to top of window |w|@>=
glog.V(debug).Infoln("go to top of the window")
w.WriteAddr("#0")
w.WriteCtl("dot=addr\nshow")

@
@<Get numbers of selected messages@>=
@<Create |msgs|@>
if (ev.Type&goacme.Tag)==goacme.Tag && len(ev.Arg) > 0 {
	s:=ev.Arg
	@<Read a message number@>
} else if err:=box.w.WriteCtl("addr=dot"); err!=nil {
	glog.Errorf("can't write to 'ctl': %s\n", err)
} else
	@<Read message numbers@>

@
@<Mark to delete messages@>=
@<Get numbers of selected messages@>
if len(msgs)!=0 {
	@<Send |msgs| to delete@>
	continue
}

@
@<Unmark messages@>=
@<Get numbers of selected messages@>
if len(msgs)!=0 {
	@<Send |msgs| to undelete@>
	continue
}

@
@<Other kinds of actions@>=
del
undel

@
@<Send |msgs| to delete@>=
glog.V(debug).Infoln("sending messages to mark for deletion")
ach<-&struct{m msgmap; a action}{msgs, del}

@
@<Send |msgs| to undelete@>=
glog.V(debug).Infoln("sending messages to unmark for deletion")
ach<-&struct{m msgmap; a action}{msgs, undel}

@
@<Handling of other types of action@>=
case del:
	var msgs messages
	for _, id:=range d.ids {
		@<Mark to delete |id| message@>
	}
	@<Refresh |msgs|@>
case undel:
	var msgs messages
	for _, id:=range d.ids {
		@<Unmark to delete |id| message@>
	}
	@<Refresh |msgs|@>

@
@<Mark to delete |id| message@>=
if p, ok:=box.all.Search(id); ok {
	if box.all[p].deleted {
		continue
	}
	box.all[p].deleted=true
	box.deleted++
	msgs=append(msgs, box.all[p])
	if box.all[p].w!=nil {
		msg:=box.all[p]
		@<Write a tag of message window@>
	}
	glog.V(debug).Infof("the '%v' message is marked for deletion\n", id)
}

@
@<Unmark to delete |id| message@>=
if p, ok:=box.all.Search(id); ok {
	if !box.all[p].deleted {
		continue
	}
	box.all[p].deleted=false
	box.deleted--
	msgs=append(msgs, box.all[p])
	if box.all[p].w!=nil {
		msg:=box.all[p]
		@<Write a tag of message window@>
	}
	glog.V(debug).Infof("the '%v' message is unmarked for deletion\n", id)

}

@ Here is processing a final deletion of messages from \.{mailfs}. Any message could be printed in
other mailboxes in threads, so we collect messages in |msgs| and send |msgs| to all mailboxes.
@<Delete messages@>=
f, err:=box.fid.Walk("ctl")
if err==nil {
	err=f.Open(plan9.OWRITE)
}
if err!=nil {
	glog.Errorf("can't open 'ctl': %v\n", err)
	continue
}
var msgs messages
for i:=0; i<len(box.all);{
	if !box.all[i].deleted || box.all[i].w!=nil {
		i++
		continue
	}
	msgs=append(msgs, box.all[i])
	@<Delete a message at position |i|@>
}
cmd:=fmt.Sprintf("delete %s", box.name)
for _, msg:=range msgs {
	cmd=fmt.Sprintf("%s %d ", cmd, msg.id)
}
glog.V(debug).Infof("command to delete messages: '%s'\n", cmd)
if _, err:=f.Write([]byte(cmd)); err!=nil{
	glog.Errorf("can't delete messages: %v\n", err)
}
f.Close()
@<Send deleted |msgs|@>


@ |mdch| is a channel receives slices of messages to delete.
@<Rest of |mailbox| members@>=
mdch	chan messages

@
@<Rest of initialization of |mailbox|@>=
mdch:make(chan messages, 100),

@ All messages from a received slice |m| will be removed from |box|'s window. In case of the thread mode
|children| is obtained and refreshed.
@<Processing of other |box| channels@>=
case m:=<-box.mdch:
	if box.w==nil {
		continue
	}
	glog.V(debug).Infof("%d messages were received to be deleted from the '%s' mailbox\n", len(m), box.name)
	for _, msg:=range m {
		@<Erase the message@>
		if box.threadMode() {
			@<Get |children| for |msg|@>
			@<Refresh |children|@>
		}
	}
	@<Check for a clean state of the |box|'s window@>

@ One message can be presented in multiple boxes, so we have to delete messages from all boxes.
|mdch| is a channel to receive signals to delete messages.
@<Variables@>=
mdch chan messages=make(chan messages, 100)

@
@<Processing of other common channels@>=
case msgs:=<-mdch:
	for i, _:=range boxes {
		glog.V(debug).Infof("sending %d messages to delete in the '%s' mailbox\n", len(msgs), boxes[i].name)
		boxes[i].mdch<-append(messages{}, msgs...)
	}

@
@<Send deleted |msgs|@>=
mdch<-msgs

@
@<Mark messages as seen@>=
@<Get numbers of selected messages@>
if len(msgs)!=0 {
	@<Send |msgs| to mark them seen@>
	continue
}

@
@<Other kinds of actions@>=
seen

@
@<Send |msgs| to mark them seen@>=
glog.V(debug).Infoln("sending messages to mark them seen")
ach<-&struct{m msgmap; a action}{msgs, seen}

@
@<Handling of other types of action@>=
case seen:
	var msgs messages
	f, err:=box.fid.Walk("ctl")
	if err==nil {
		err=f.Open(plan9.OWRITE)
	}
	if err!=nil {
		glog.Errorf("can't open 'ctl' file of the '%s' messagebox: %v\n", box.name, err)
		continue
	}
	for _, id:=range d.ids {
		p, ok:=box.all.Search(id)
		if !ok || !box.all[p].unread {
			continue
		}
		msgs=append(msgs, box.all[p])
	}
	cmd:="read"
	for _, v:=range msgs {
		cmd+=fmt.Sprintf(" %d", v.id)
	}
	if _, err:=f.Write([]byte(cmd)); err!=nil {
		glog.Errorf("can't write to 'ctl' file of the '%s' messagebox: %v\n", box.name, err)
	}
	for _, msg:=range msgs {
		id:=msg.id
		@<Remove |id| message from |unread|@>
		@<Refresh the message's view@>
	}
	f.Close()
	@<Send |box| to refresh the main window@>
	@<Refresh |msgs|@>

@* Linking of threads.

Here we define global map of unique message identifiers on a pointer to messages, its parents and children.
An unique id of every message will be stored in this map.
It will be changed in a separated goroutine, so a corresponding channel |idch| is defined too.
|idch| operates with a pair contains a message |msg| and |val interface{}|. We can use |val| with different
types of data to do different operations with |idmap|.
In case of an addition of a new message the |val| is set to |chan bool|. The channel is used to inform the sender
the message has been added successfully.
In case of a deletion of the message |val| is set to |nil|.
To get children of the message we set |val| to a channel of |idmessages|.
To get a parent of the message we set |val| to a channel of |parentmsg|.
To get a root of a thread with the message we set |val| to a channel of |rootmsg|.
To get a level of the message in the thread we set |val| to a channel of |int|.

@
@<Types@>=
idmessages	[]*message
rootmsg		*message
parentmsg	*message

@ |idlinks| struct contains |msgs| slice of messages with the same |messageid|, a pointer to |idlinks| of parent and
|children| slice of pointers to |idlinks| of children.
@<Types@>=
idlinks struct{
	msgs messages
	parent *idlinks
	children []*idlinks
}

@
@<Variables@>=
idmap=make(map[string]*idlinks)
idch=make(chan struct{msg *message; val interface{}}, 100)

@
@<Rest of |message| members@>=
inreplyto 	string
messageid 	string

@
@<Imports@>=
"errors"

@ Here we read |inreplyto| and |messageid|. When |messageid| is obtained, we send |msg| to |idch|.
We assume |inreplyto| is already filled, because it is placed above of |messageid| in the |"info"| file.
@<Read other fields of a message@>=
{
	if _, err:=fmt.Sscanf(s, "inreplyto %s", &msg.inreplyto); err==nil {
		msg.inreplyto=strings.Trim(msg.inreplyto, "<>")
		continue
	}
	if _, err:=fmt.Sscanf(s, "messageid %s", &msg.messageid); err==nil {
		msg.messageid=strings.Trim(msg.messageid, "<>")
		ch:=make(chan bool)
		idch<-struct{msg *message; val interface{}}{msg, ch}
		if ok:=<-ch; !ok {
			return nil, false, errors.New(fmt.Sprintf("a message '%s' is duplicated", msg.messageid))
		}
		continue
	}
}

@ Processing of |idch| in a separated goroutine..
@<Start a collector of message identifiers@>=
go func() {
	for {
		select {
			@<On exit?@>
			case v:=<-idch:
				if v.val==nil {
					@<Clean an entry with |v.msg.messageid| from |idmap|@>
				} else if ch, ok:=v.val.(chan bool); ok {
					@<Append a message with |v.msg.messageid| to |idmap|@>
				} else if ch, ok:=v.val.(chan idmessages); ok {
					@<Send |children|@>
				} else if ch, ok:=v.val.(chan parentmsg); ok {
					@<Send |parent|@>
				} else if ch, ok:=v.val.(chan rootmsg); ok {
					@<Send |root|@>
				} else if ch, ok:=v.val.(chan int); ok {
					@<Send |level|@>
				}
		}
	}
}()

@ When |msg| is appended we should check if |v.id| already exists. It can exist if there are
duplicated messages or there are children for this |v.id|. For the latter case an entry is added
to |idmap| with empty |msgs|. Later, when a message with |v.id| is added, we just add the new |msg| in |msgs|.

If |msg| has |inreplyto| is filled, we add |val| pointer to |children| of |msg.inreplyto| message and
set a parent for |val|.

@<Append a message with |v.msg.messageid| to |idmap|@>=
{
	glog.V(debug).Infof("appending a '%s' ('%s/%d') message to idmap\n", v.msg.messageid, v.msg.box.name, v.msg.id)
	val, ok:=idmap[v.msg.messageid]
	if !ok {
		glog.V(debug).Infof("'%s' ('%s/%d') message  doesn't exist, creating\n", v.msg.messageid, v.msg.box.name, v.msg.id)
		val=new(idlinks)
		idmap[v.msg.messageid]=val
	}
	if len(val.msgs)>0 {
		glog.V(debug).Infof("%v(%v) is a duplicate of %v(%v)\n",
			v.msg.id, v.msg.messageid,
			val.msgs[0].id, val.msgs[0].messageid)
	}
	val.msgs=append(val.msgs, v.msg)

	if len(v.msg.inreplyto)!=0 && len(val.msgs)==1 {
		pval, ok:=idmap[v.msg.inreplyto]
		if !ok {
			pval=new(idlinks)
			idmap[v.msg.inreplyto]=pval
		}
		
		pval.children=append(pval.children, val)
		val.parent=pval
	}
	ch<-true
}

@ When we are removing a message, we just remove |msg| from |msgs| slice.
But if |msgs| slice is empty, we have to clean pointer to the entry in a parent and chindren entries.
We leave a non-empty entry in |idmap| to store links of children.
@<Clean an entry with |v.msg.messageid| from |idmap|@>=
{
	val, ok:=idmap[v.msg.messageid]
	if !ok {
		continue
	}
	for i, _:=range val.msgs {
		if val.msgs[i]==v.msg {
			val.msgs.Delete(i)
			break
		}
	}
	if len(val.msgs)>0 {
		continue
	}
	if val.parent!=nil {
		for i, _:=range val.parent.children {
			if val.parent.children[i]==val {
				val.parent.children=append(val.parent.children[:i], val.parent.children[i+1:]...)
				break
			}
		}
	}
	for _, ch:=range val.children {
		ch.parent=nil
	}
	if len(val.children)==0 {
		delete(idmap, v.msg.messageid)
	}
}

@ A few methods have to be implemented for |idmessages| to have an ability to sort of them in order of increasing of date.
@c
func (this idmessages) Len() int {
	return len(this)
} @#

func (this idmessages) Less(i, j int) bool {
	return this[i].date.Unix() < this[j].date.Unix()
} @#

func (this idmessages) Swap(i, j int) {
	t:=this[i]
	this[i]=this[j]
	this[j]=t
} @#


@ If there is |v.msg.messageid| in |idmap|, we fill |children| with correspinding children and sort them in order of increasing of date
@<Send |children|@>=
{
	if m, ok:=idmap[v.msg.messageid]; ok {
		var children idmessages
		for _, val:=range m.children {
			@<Get |mgs| with the same box like |v.msg.box| or the first one@>
			children=append(children, msg)		
		}
		sort.Sort(children)
		glog.V(debug).Infof("sending %d children for '%s'\n", len(children), v.msg.messageid)
		ch<-children
	} else {
		glog.V(debug).Infof("'%s' is not found\n", v.msg.messageid)
		ch<-nil
	}
}

@
@<Get |children| for |msg|@>=
var children idmessages
{
	ch:=make(chan idmessages)
	glog.V(debug).Infof("getting children for '%s'\n", msg.messageid)
	idch<-struct{msg *message; val interface{}}{msg, ch}
	children=<-ch
}

@
@<Send |parent|@>=
{
	val, ok:=idmap[v.msg.messageid]
	if !ok {
		glog.V(debug).Infof("'%s' is not found\n", v.msg.messageid)
		ch<-nil
	}
	if val.parent==nil || len(val.parent.msgs)==0 {
		glog.V(debug).Infof("'%s' hasn't got a parent\n", v.msg.messageid)
		ch<-nil
	} else {
		val=val.parent
		@<Get |mgs| with the same box like |v.msg.box| or the first one@>
		glog.V(debug).Infof("sending parent '%s' for '%s'\n", msg.messageid, v.msg.messageid)
		ch<-msg
	}
}

@ We are looking for a message with the same |box|.
@<Get |mgs| with the same box like |v.msg.box| or the first one@>=
var msg *message
if val!=nil && len(val.msgs)>0 {
	msg=val.msgs[0]
	for i, _:=range val.msgs {
		if val.msgs[i].box==v.msg.box {
			msg=val.msgs[i]
			break
		}
	}
}

@
@<Get |parent| for |msg|@>=
var parent *message
{
	ch:=make(chan parentmsg)
	glog.V(debug).Infof("getting parent for '%s'\n", msg.messageid)
	idch<-struct{msg *message; val interface{}}{msg, ch}
	parent=<-ch
}

@
@<Send |root|@>=
{
	if val, ok:=idmap[v.msg.messageid]; ok {
		for val.parent!=nil && len(val.parent.msgs)>0 {
			val=val.parent
		}
		@<Get |mgs| with the same box like |v.msg.box| or the first one@>
		if msg==nil {
			msg=v.msg
		}
		ch<-rootmsg(msg)
	} else {
		glog.V(debug).Infof("'%s' is not found\n", v.msg.messageid)
		ch<-nil
	}
}

@
@<Get |root| of |msg|@>=
root:=msg
{
	ch:=make(chan rootmsg)
	glog.V(debug).Infof("getting root for '%s' ('%s/%d')\n", msg.messageid, msg.box.name, msg.id)
	idch<-struct{msg *message; val interface{}}{msg, ch}
	root=<-ch
}

@
@<Send |level|@>=
{
	if val, ok:=idmap[v.msg.messageid]; ok {
		level:=0
		for val.parent!=nil && len(val.parent.msgs)>0 {
			val=val.parent
			level++
		}
		glog.V(debug).Infof("sending level '%d' for '%s' ('%s/%d')\n", level, v.msg.messageid, v.msg.box.name, v.msg.id)
		ch<-level
	} else {
		glog.V(debug).Infof("'%s' is not found\n", v.msg.messageid)
		ch<-0
	}
}

@
@<Get |level| for |msg|@>=
var level int
{
	ch:=make(chan int)
	glog.V(debug).Infof("getting root for '%s'\n", msg.messageid)
	idch<-struct{msg *message; val interface{}}{msg, ch}
	level=<-ch
}

@ Here we send |msg| to |idch| with |nil| like a |val| to clean up a thread links.
@<Clean up |msg|@>=
glog.V(debug).Infof("cleaning up the '%d' message\n", msg.id)
if msg!=nil {
	idch<-struct{msg *message; val interface{}}{msg: msg}
}

@* Printing of messages.

Printing of the messages is a kind of trick. To avoid of locks of |box|'s stuff the print is produced in the |box|'s
message loop. |rfch| is a channel receives slice of messages have to be printed and flag to seek a position
to start a print or to print in the end.
A data from |rfch| is redirected to an internal channel |irfch|.
A slice of messages is sent to |rfch|, the |box|'s message loop reads the slice and print at most |100| messages,
then resend the rest to |rfch|. If we need to stop printing of messages, we drop the rest
of a printing queue by recreation of |irfch|.

@ |refresh| holds flags point out how to print |msgs|
@<Types@>=
refreshFlags int

refresh	struct {
	flags refreshFlags
	msgs messages
}

@ |seek| means a position of the message should be determinated,
|insert| means the message should be inserted if the position is not found,
|exact| means the message should be inserted only if its exact position is found.
@<Constants@>=
const (
	seek refreshFlags=1<<iota
	insert refreshFlags=1<<iota
	exact refreshFlags=1<<iota
)

@
@<Rest of |mailbox| members@>=
rfch	chan *refresh
irfch	chan *refresh
reset	bool

@
@<Rest of initialization of |mailbox|@>=
rfch:make(chan *refresh, 100), @/
irfch:make(chan *refresh, 100),

@ |box.rfch| receives a slice of messages to be printed.
In case of threaded messages should be printed, but linking of messages
still hasn't finished, the slice is ignored. Actually |box.rfch| is an external
channel, it resend a data into |box.irfch|. If we need to stop printing,
we just recreate |box.irfch|.

@<Processing of other |box| channels@>=
case v:=<-box.rfch:
	box.irfch<-v

case v:=<-box.irfch:
	glog.V(debug).Infof("a signal to print message of the '%s' mailbox window has been received\n", box.name)
	if box.w==nil {
		glog.V(debug).Infof("a window of the '%s' mailbox doesn't exist, ignore the signal\n", box.name)
		continue
	}
	if	box.threadMode() && !counted {
		glog.V(debug).Infof("counting of threads of the '%s' mailbox is not finished, ignore the signal\n", box.name)
		continue
	}
	@<Print messages from |v.msgs|@>


@
@<Determine of |src|@>=
var src messages
if box.shownew {
	src=box.unread
} else {
	src=box.all
}

@ |messages.Check| checks for message with |messageid| is already exists in |this|
If it doesn't exist it is added and |true| is returned, |false| otherwise.
@c
func (this *messages) Check(msg *message) bool {
	pos:=sort.Search(len(*this), func(i int) bool {return (*this)[i].messageid<=msg.messageid});
	if pos!=len(*this) && (*this)[pos].messageid==msg.messageid {
		return false
	}
	*this=append(*this, nil)
	copy((*this)[pos+1:], (*this)[pos:])
	(*this)[pos]=msg
	return true
}


@ All enumerated messages should be printed according to the options.
In case of the thread mode sequences of full threads should be made.
To avoid of duplicating of threads roots of threads are accumulated and every
new root is checked for an existance.
@<Inform |box| to print messages@>=
{
	glog.V(debug).Infof("inform the '%s' mailbox to print messages\n", box.name)
	@<Determine of |src|@>
	msgs:=append(messages{}, src...)
	if box.threadMode() {
		src=msgs
		msgs=make(messages, 0, len(src))
		var roots messages
		for len(src)>0 {
			msg:=src[0]
			@<Get |root| of |msg|@>
			if !roots.Check(root) {
				@<Remove |msg| from |src|@>
				continue
			}
			glog.V(debug).Infof("root of thread: '%s/%d'\n", root.box.name, root.id)
			@<Make a full thread in |msgs| for |root|@>
		}
	} else {
		@<Set |pos| of |box|@>
	}
	box.rfch<-&refresh{0, msgs}
}

@ |msg| is added to the |msgs| list and all its children are processed.
To avoid duplicates |msg| has to be removed from |src|.
@<Make a full thread in |msgs| for |root|@>=
msgs=append(msgs, root)
@<Remove |msg| from |src|@>
msgs, src=getchildren(root, msgs, src)

@
@<Remove |msg| from |src|@>=
if p, ok:=src.Search(msg.id); ok {
	glog.V(debug).Infof("removing '%d' from src\n", src[p].id)
	src.Delete(p)
}

@ |getchildren| gets children for |msg|, removes |msg| from |src| and does the same
for all children recursively.
@c
func getchildren(msg *message, dst messages, src messages) (messages, messages) {
	@<Get |children| for |msg|@>
	for _, msg:=range children {
		dst=append(dst, msg)
		@<Remove |msg| from |src|@>
		dst, src=getchildren(msg, dst, src)
	}
	return dst, src
}

@ A list with full thread of messages is made for |msg|
@<Inform |box| to print a full thread with |msg|@>=
@<Get |root| of |msg|@>
var msgs messages
src:=append(messages{}, root)
@<Make a full thread in |msgs| for |root|@>
box.rfch<-&refresh{0, msgs}

@ Only |msg| should be printed.
@<Inform |box| to print |msg|@>=
{
	glog.V(debug).Infof("inform the '%s' mailbox to print a message '%d'\n", box.name, msg.id)
	box.rfch<-&refresh{seek|insert, append(messages{}, msg)}
}

@
@<Inform |box| to print |msgs|@>=
{
	glog.V(debug).Infof("inform the '%s' mailbox to print messages\n", box.name)
	box.rfch<-&refresh{seek|insert, msgs}
}

@ Only |msg| should be refreshed.
@<Refresh |msg|@>=
{
	glog.V(debug).Infof("refresh a message '%d'\n",msg.id)
	mrfch<-&refresh{seek, append(messages{}, msg)}
}

@ |msgs| will be refreshed in |box| window with setting a position for every message if is found.
@<Inform |box| to refresh |msgs|@>=
{
	if len(msgs)!=0 {
		glog.V(debug).Infof("inform the '%s' mailbox to refresh messages\n", box.name)
		box.rfch<-&refresh{seek, msgs}
	}
}

@ |msgs| will be refresh with setting a position for every message if is found.
@<Refresh |msgs|@>=
{
	if len(msgs)!=0 {
		glog.V(debug).Infoln("refresh messages\n")
		mrfch<-&refresh{seek, msgs}
	}
}

@ |msg| will be printed only if the exact position is found.
@<Print |msg| at exact positon@>=
{
	glog.V(debug).Infof("print '%s/%d' at exact position\n", box.name, msg.id)
	mrfch<-&refresh{seek|insert|exact, append(messages{}, msg)}
}

@ One message can be presented in multiple boxes, so we have to refresh messages in all boxes.
|mrfch| is a channel to receive signals to refresh messages.
@<Variables@>=
mrfch chan *refresh=make(chan *refresh)

@
@<Processing of other common channels@>=
case r:=<-mrfch:
	for i, _:=range boxes {
		glog.V(debug).Infof("sending messages to refresh in the '%s' mailbox\n", boxes[i].name)
		boxes[i].rfch<-&refresh{r.flags, append(messages{}, r.msgs...)}
	}



@ We need to store a current position of |src| to know a message will be started to print with.
@<Rest of |mailbox| members@>=
pos	int

@
@<Clean window-specific stuff@>=
box.pos=0
ontop=false

@ We set |pos| to len of determinated |src| to avoid of printing messages twice.
@<Set |pos| of |box|@>=
box.pos=len(src)

@ Printing during the counting process is made only for plain mode. We use |box.pos| like a position of a first
message to print and print a number of messages is multiple of |500|
@<Inform |box| to print counting messages@>=
if !box.threadMode() {
	@<Determine of |src|@>
	if len(src)!=0 && box.pos<len(src) && len(src)%500==0  {
		glog.V(debug).Infof("inform the '%s' mailbox to print the last %d messages\n", box.name, len(src)-box.pos)
		msgs:=append(messages{}, src[box.pos:len(src)]...)
		@<Set |pos| of |box|@>
		box.rfch<-&refresh{0, msgs}
	}
}

@ Here we print the rest of counted messages.
@<Inform |box| to print the rest of counting messages@>=
if !box.threadMode() {
	@<Determine of |src|@>
	if box.pos<len(src) {
		glog.V(debug).Infof("inform the '%s' mailbox to print the last %d messages\n", box.name, len(src)-box.pos)
		msgs:=append(messages{}, src[box.pos:len(src)]...)
		@<Set |pos| of |box|@>
		box.rfch<-&refresh{0, msgs}
	}
}

@
@<Constants@>=
const eof="$"

@
@<Print messages from |v.msgs|@>=
{
	glog.V(debug).Infof("printing of messages of the '%s' mailbox from v.msgs, len(v.msgs): %d, with flags: %v\n", box.name, len(v.msgs), v.flags)
	f, err:=box.w.File("data")
	if err!=nil {
		glog.Errorf("can't open 'data' file of the '%s' messagebox: %v\n", box.name, err)
		continue
	}
	if (v.flags&seek)==seek {
		@<Write a tag of |box| window@>
		msg:=v.msgs[0]
		@<Trying to find a place for |msg| in the |box| window@>
	} else if err:=box.w.WriteAddr(eof); err!=nil {
		glog.Errorf("can't write '%s' to 'addr' file: %s\n", eof, err)
		continue
	}
	w:=box.w
	glog.V(debug).Infof("printing of messages of the '%s' mailbox\n", box.name)
	buf:=make([]byte, 0, 0x8000)
	@<Compose messages of the |box|@>
	if _, err:=f.Write(buf); err!=nil {
		glog.Errorf("can't write to 'data' file of the '%s' messagebox: %v\n", box.name, err)
	}
	@<Go to the top of window for first 100 messages@>
	@<Send a rest of |msgs|@>
}

@
@<Send a rest of |msgs|@>=
if len(v.msgs)>0 {
	box.rfch<-&refresh{v.flags, v.msgs}
} else {
	@<Check for a clean state of the |box|'s window@>
}

@ To stay on top of the box's window when printing we go to top for first
100 messages, I hope it is enough to print other messages in the background without scrolling.
@<Go to the top of window for first 100 messages@>=
if !ontop {
	glog.V(debug).Infof("pcount: %v, ontop: %v\n", pcount, ontop)
	@<Go to top of window |w|@>
	if pcount>=100 {
		ontop=true
	}
}


@ Here the messages composing is produced. To avoid of overloading of events processing
we print a lot of messages at a time. But if |seek| is set in |v.flags| messages will be printed one
by one, because we have to set a position for every message..

@<Compose messages of the |box|@>=
c:=0
for len(v.msgs)>0 && c<100 {
	msg:=v.msgs[0]
	glog.V(debug).Infof("printing of '%s/%d' message with in the '%s' mailbox\n", msg.box.name, msg.id, box.name)
	if box.threadMode() {
		@<Add the thread level marks@>
	}
	c++
	@<Compose a header of |msg|@>
	v.msgs=v.msgs[1:]
	if (v.flags&seek)==seek {
		break
	}
}
pcount+=c

@
@<Clean window-specific stuff@>=
{
	glog.V(debug).Infof("clean window-specific stuff of the '%s' mailbox\n", box.name)
	close(box.irfch)
	box.irfch=make(chan *refresh, 100)
	pcount=0
	ontop=false
}

@
@<Add the thread level marks@>=
{
	@<Get |level| for |msg|@>
	for ;level>0; level--{
		buf=append(buf, levelmark...)
	}
}

@ In case deleted message has children we should refresh views of these children.
So we compose a list of messages and send them to refresh. But if a child is not belonged to |box|
we have to erase it instead of refreshing.
@<Refresh |children|@>=
{
	if len(children)!=0 {
		var msgs messages
		var src messages
		for _, msg:=range children {
			if msg.box!=box {
				@<Erase the message@>
			} else {
				root:=msg
				@<Make a full thread in |msgs| for |root|@>
			}
		}
		@<Inform |box| to refresh |msgs|@>
	}
}

@
@<Write a tag of main window@>=
glog.V(debug).Infoln("writing a tag of the main window")
if err:=writeTag(mw, fmt.Sprintf(" %s %s ", @t\1@>@/
	func() string {if shownew {return "ShowAll"} else {return "ShowNew"}}(), @/
	func() string {if showthreads {return "ShowPlain"} else {return "ShowThreads"}}()) @t\2@>);
	err!=nil {
	glog.Errorf("can't set a tag of the main window: %v", err)
}

@
@c
func writeTag(w *goacme.Window, t string) error {
	if w==nil {
		return nil
	}
	tag, err:=w.File("tag")
	if err!=nil {
		return err
	}
	if err:=w.WriteCtl("cleartag"); err!=nil {
		return err
	}
	_, err=tag.Write([]byte(t))
	return err
}

@ |deleted| contains a count of messages to delete.
@<Rest of |mailbox| members@>=
deleted int

@
@<Write a name of |box| window@>=
name:="Amail/"+box.name
w:=box.w
@<Print the |name| for window |w|@>

@
@<Write a tag of |box| window@>=
box.writeTag(counted)

@
@c
func (box *mailbox) writeTag(counted bool) {
	glog.V(debug).Infof("write a tag of the '%s' mailbox's window\n", box.name)
	@<Determine of |src|@>
	if err:=writeTag(box.w, fmt.Sprintf(" %sMail %s%s%s%s%s%sSearch ", @t\1@>@/
		func() string {
			if box.deleted>0 {
				return "Put "
			}
			return ""
		}(), @/
		func() string {
			if len(src)>0 {
				return "Delmesg "
			}
			return ""
		}(), @/
		func() string {
			if len(src)>0 && box.deleted>0 {
				return "UnDelmesg "
			}
			return ""
		}(), @/
		func() string {
			if box.thread {
				if box.shownew {
					return "ShowNew "
				} else {
					return "ShowAll "
				}
			} else if len(src)>0 && counted && (box.shownew || !box.showthreads) {
				return "Thread "
			}
			return ""
		}(), @/
		func() string {
			if box.showthreads && !counted {
				return ""
			}
			if box.shownew {
				return "ShowAll "
			} else {
				return "ShowNew "
			}
		}(), @/
		func() string {
			if box.showthreads {
				return "ShowPlain "
			} else if counted {
				return "ShowThreads "
			} else {
				return ""
			}
		}(), @/
		func() string {
			if len(src)>0 {
				return "Seen "
			}
			return ""
		}()) @t\2@>)
		err!=nil {
		glog.Errorf("can't set a tag of the '%s' box's window: %v\n", box.name, err)
	}
}

@
@<Clean |box| window@>=
glog.V(debug).Infof("clean the '%s' mailbox's window\n", box.name)
clean(box.w)

@
@c
func clean(w *goacme.Window){
	if err:=w.WriteAddr(wholefile); err!=nil {
		glog.Errorf("can't write '%s' to 'addr' file: %s\n", wholefile, err)
	} else if data, err:=w.File("data"); err!=nil {
		glog.Errorf("can't open 'data' file: %s\n", err)
	} else if _, err:=data.Write([]byte("")); err!=nil {
		glog.Errorf("can't write to 'data' file: %s\n", err)
	}
}

@
@<Constants@>=
const bof="#0-"
const eol="+#0"

@ For the first we try to find the message itself. If the message is new and |insert| is set in |v.flags|, we should
find its neighbours and set address according to the position. The message should be skipped in other mailboxes
if it is not the thread mode.
@<Trying to find a place for |msg| in the |box| window@>=
@<Determine of |src|@>
@<Compose |addr|@>
glog.V(debug).Infof("composed message addr '%s' in the '%s' mailbox\n", addr, box.name)
if err:=box.w.WriteAddr(addr); err!=nil {
	glog.V(debug).Infof("the '%d' message is not found in the window of the '%s' mailbox\n", msg.id, box.name)
	if (v.flags&insert)==0 {
		@<Skip current message@>
	}
	if box.threadMode() {
		@<Set a position for a threaded message@>
	} else if msg.box!=box {
		@<Skip current message@>
	} else if p, ok:=src.Search(msg.id); !ok {
		glog.V(debug).Infof("the '%d' message is not found in  the '%s' mailbox's window\n", msg.id, box.name)
	} else if p==0 {
		if err:=box.w.WriteAddr(bof); err!=nil {
			glog.Errorf("can't write '%s' to 'addr' file of the '%s' mailbox's window: %s\n", bof, box.name,err)
		}
	} else if p==len(src)-1 {
		if err:=box.w.WriteAddr(eof); err!=nil {
			glog.Errorf("can't write '%s' to 'addr' file of the '%s' mailbox's window: %s\n", eof, box.name, err)
		}
	} else {
		msg:=src[p-1]
		@<Compose |addr|@>
		addr+=eol
		if err:=box.w.WriteAddr(addr); err!=nil {
			glog.V(debug).Infof("can't write '%s' to 'addr' of the '%s' mailbox's window: %v\n", addr, box.name, err)
		}
	}
}

@
@<Skip current message@>=
glog.V(debug).Infof("the '%d' message won't be inserted in the '%s' mailbox's window\n", v.msgs[0].id, box.name)
v.msgs=v.msgs[1:]
@<Send a rest of |msgs|@>
continue


@
@<Compose |addr|@>=
addr:=fmt.Sprintf("0/^[%s]*(%s)?%s%d(%s)?\\/.*\\n.*\\n/",  @t\1@>@/
			escape(levelmark), @/
			escape(deleted), @/
			func() string { if box!=msg.box {return escape(msg.box.name+"/")}; return ""}(), @/
			msg.id, @/
			escape(newmark)@t\2@>)

@ If |msg| has a parent, it should be printed after last child of the thread.
In case of |msg| is only child of |msg.parent|, |msg| will be printed after |msg.parent|.
If |msg| has no parent and |exact| is not set in |v.flags| it will be printed on top of the window.
@<Set a position for a threaded message@>=
@<Get |parent| for |msg|@>
if parent!=nil {
	glog.V(debug).Infof("msg '%d' has a parent, looking for a position to print\n", msg.id)
	m:=msg
	msg=parent
	found:=false
	for !found {
		@<Get |children| for |msg|@>
		if len(children)==0 {
			break
		}
		for i, v:=range children {
			if v==m {
				if i==0 {
					found=true
				}
				break
			}
			msg=v
		}
	}
	glog.V(debug).Infof("the '%d' message will be printed after the '%d' message\n", m.id, msg.id)
	@<Compose |addr|@>
	addr+=eol
	if err:=box.w.WriteAddr(addr); err!=nil {
		glog.V(debug).Infof("can't write '%s' to 'addr' of the '%s' mailbox's window: %v\n", addr, box.name, err)
		if (v.flags&exact)==exact {
			@<Skip current message@>
		}
	}
} else if (v.flags&exact)==exact {
	@<Skip current message@>
} else if err:=box.w.WriteAddr(bof); err!=nil {
	glog.Errorf("can't write '%s' to 'addr' file of the '%s' mailbox's window: %v\n", bof, box.name, err)
}

@
@<Search messages@>=
{
	msgs:=box.search(ev.Arg)
	@<Clean |box| window@>
	@<Clean window-specific stuff@>
	name:=fmt.Sprintf("Amail/%s/Search(%s)", box.name, strings.Replace(ev.Arg, " ", "␣", -1))
	w:=box.w
	box.thread=false
	box.shownew=true
	box.showthreads=false
	@<Print the |name| for window |w|@>
	glog.V(debug).Infof("len of msgs: %v\n", len(msgs))
	box.rfch<-&refresh{0, msgs}
}

@
@c
func (box *mailbox) search(str string) (msgs messages) {
	if len(str)==0 {
		return
	}
	f, err:=box.fid.Walk("search")
	if err==nil {
		err=f.Open(plan9.ORDWR)
	}
	if err!=nil {
		glog.Errorf("can't open 'search' file: %s\n", err)
		return
	}
	defer f.Close()
	if _, err:=f.Write([]byte(str)); err!=nil {
		glog.Errorf("can't write to 'search' file: %s\n", err)
	}
	b:=bufio.NewReader(f)
	for s, err:=b.ReadString(' '); err==nil || err==io.EOF; s, err=b.ReadString(' ') {
		s=strings.TrimSpace(s)
		glog.V(debug).Infoln("search: ", s)
		if num, err:=strconv.Atoi(s); err==nil {
			if p, ok:=box.all.Search(num); ok {
				msgs.Insert(box.all[p], 0)
			}
		}
		if err==io.EOF {
			break
		}
	}
	return
}

@* Viewing of a message.

At first let's extend |mailbox| by a |lch| channel
@<Rest of |mailbox| members@>=
ach chan *struct{ids []int; a action}

@
@<Rest of initialization of |mailbox|@>=
ach:make(chan *struct{ids []int; a action}, 100),

@ We have to extend |message| too by |*goacme.Window|
@<Rest of |message| members@>=
w *goacme.Window

@ Here we will process requests to open messages. If the message is new, it should be removed from |box.unread| and
its view in all windows should be changed. The count of unread messages on the main window should be refreshed too.
We accumulate messages with changed status in |msgs| and refresh them after all messages are opened.
@<Processing of other |box| channels@>=
case d:=<-box.ach:
	switch d.a {
		case view:
			var msgs messages
			for _, id:=range d.ids {
				glog.V(debug).Infof("opening a window with the '%d' message of the '%s' mailbox\n", id, box.name)
				p, ok:=box.all.Search(id)
				if !ok {
					glog.V(debug).Infof("the '%d' message of the '%s' mailbox has not found\n", id, box.name)
					continue
				}
				msg:=box.all[p]
				if msg.w==nil {
					if msg.unread {
						@<Remove |id| message from |unread|@>
						@<Refresh the message's view@>
						@<Send |box| to refresh the main window@>
					}
					if err:=msg.open(); err!=nil {
						continue
					}
				} else {
					glog.V(debug).Infof("a window of the '%d' message of the '%s' already exists, just show it\n", id, box.name)
					msg.w.WriteCtl("dot=addr\nshow")
				}
			}
			@<Refresh |msgs|@>
		@<Handling of other types of action@>
	}

@
@<Remove |id| message from |unread|@>=
msg.unread=false
box.unread.DeleteById(id)

@ In case of viewing new messages only we have to remove the message from window.
Also |msg| has to be added to |msgs| to refresh the message's view in other windows.
@<Refresh the message's view@>=
if !box.thread && box.shownew {
	@<Erase the message@>
	@<Check for a clean state of the |box|'s window@>
}
msgs=append(msgs, msg)



@
@<Check for a clean state of the |box|'s window@>=
{
	glog.V(debug).Infof("box.deleted:%d\n", box.deleted)
	@<Write a tag of |box| window@>
	w:=box.w
	if box.deleted==0 {
		@<Set window |w| to clean state@>
	} else {
		@<Set window |w| to dirty state@>
	}
}

@ Here we remove a message |msg| from |box|'s window.
@<Erase the message@>=
box.eraseMessage(msg)


@
@c
func (box *mailbox) eraseMessage(msg *message){
	if box.w==nil {
		return
	}
	glog.V(debug).Infof("removing the '%d' message of the '%s' mailbox from the '%s' mailbox\n", @t\1@>@/
		msg.id, msg.box.name, box.name @t\2@>)
	@<Compose |addr|@>
	if err:=box.w.WriteAddr(addr); err!=nil {
		glog.V(debug).Infof("can't write '%s' to 'addr' of the '%s' mailbox's window: %v\n", addr, box.name, err)
	} else if data, err:=box.w.File("data"); err !=nil {
		glog.Errorf("can't open 'data' file of the box '%s': %s\n", box.name, err)
	} else if _, err:=data.Write([]byte{}); err!=nil {
		glog.Errorf("can't write to 'data' file of the box '%s': %s\n", box.name, err)
	}
}

@
@<Rest of |message| members@>=
to		[]string
cc		[]string

@
@<Read other fields of a message@>=
if strings.HasPrefix(s, "to ") {
	msg.to=split(s[len("to "):])
	continue
}
if strings.HasPrefix(s, "cc ") {
	msg.cc=split(s[len("cc "):])
	continue
}


@ |split| splits |s| to a |[]string| of mail addresses that can contain a name and an address.
If a name is just |''|, it is removed.
@c
func split(s string) (strs []string) {
	f:=strings.Fields(s)
	m:=""
	for _, v:=range f {
		if strings.Contains(v, "@@") {
			m+=v
			strs=append(strs, m)
			m=""
		} else if v!="''" {
			m+=v+" "
		}
	}
	return
}

@ |open| opens a message in a separated window.
@c
func (msg *message) open() (err error) {
	glog.V(debug).Infof("open: trying to open '%d' directory\n", msg.id)
	bfid, err:=msg.box.fid.Walk(fmt.Sprintf("%d", msg.id))
	if err!=nil {
		glog.Errorf("can't walk to '%s/%d': %v\n", msg.box.name, msg.id, err)
		return err
	}
	defer bfid.Close()
	isnew:=msg.w==nil
	if isnew {
		if msg.w, err=goacme.New(); err!=nil {
			glog.Errorf("can't create a window: %v\n", err)
			return err
		}
	} else {
		@<Clean |msg.w| window@>
	}
	buf:=make([]byte, 0, 0x8000)
	@<Compose a header of the message@>
	@<Compose a body of the message@>
	w:=msg.w
	name:=fmt.Sprintf("Amail/%s/%d", msg.box.name, msg.id)
	@<Print the |name| for window |w|@>
	@<Write a tag of message window@>
	w.Write(buf)
	@<Set window |w| to clean state@>
	@<Go to top of window |w|@>
	if isnew {
		@<Start a goroutine to process events from the message's window@>
	}
	return
}

@
@<Write a tag of message window@>=
msg.writeTag()

@
@<Get a previous |pmsg|@>=
pmsg:=msg.prev()

@
@c
func (this *message) prev() (pmsg *message) {
	msg:=this
	@<Get |parent| for |msg|@>
	if parent==nil {
		return
	}
	msg=parent
	@<Get |children| for |msg|@>
	for _, v:=range children {
		if v==this {
			break
		}
		pmsg=v
	}
	return
}

@
@<Get a next |nmsg|@>=
nmsg:=msg.next()

@
@c
func (this *message) next() (nmsg *message) {
	msg:=this
	@<Get |parent| for |msg|@>
	if parent==nil {
		return
	}
	msg=parent
	@<Get |children| for |msg|@>
	for i:=0; i<len(children); i++ {
		if children[i]!=this {
			continue
		}
		i++
		if i<len(children) {
			nmsg=children[i]
		}
		break
	}
	return
}

@
@c
func (msg *message) writeTag() {
	glog.V(debug).Infof("writing a tag of the '%d' message's window\n", msg.id)
	if err:=writeTag(msg.w, fmt.Sprintf(" Q Reply all %s %s%s%s%s%s%sSave ", @t\1@>@/
		func() string {if msg.deleted {return "UnDelmesg"} else {return "Delmesg"}}(),
		func() string {
			if len(msg.text)==0 || len(msg.html)==0 {
				return ""
			} else if msg.showhtml {
				return "Text "
			} else {
				return "Html "
			}
		}(), @/
		func() string {if len(msg.html)!=0 {return "Browser "}; return ""}(), @/
		func() string {
			@<Get |parent| for |msg|@>
			if parent!=nil {
				return "Up "
			}
			return ""
		}(), @/
		func() string {
			@<Get |children| for |msg|@>
			if len(children)!=0 {
				return "Down "
			}
			return ""
		}(), @/
		func() string {
			@<Get a previous |pmsg|@>
			if pmsg!=nil {
				return "Prev "
			}
			return ""
		}(), @/
		func() string {
			@<Get a next |nmsg|@>
			if nmsg!=nil {
				return "Next "
			}
			return ""
		}() @t\2@>));
		err!=nil {
		glog.Errorf("can't set a tag of the message window: %v", err)
	}
}

@
@<Compose a header of the message@>=
{
	glog.V(debug).Infof("composing a header of the '%d' message\n", msg.id)
	buf=append(buf, fmt.Sprintf("From: %s\nDate: %s\nTo: %s\n%sSubject: %s\n\n\n", @t\1@>@/
		msg.from, msg.date, strings.Join(msg.to, ", "), @/
		func() string{if len(msg.cc)!=0{return fmt.Sprintf("CC: %s\n", strings.Join(msg.cc, ", "))};return ""}(), @/
		msg.subject)...@t\2@>)
}

@
@<Start a goroutine to process events from the message's window@>=
go func() {
	glog.V(debug).Infof("starting a goroutine to process events from the '%d' message's window\n", msg.id)
	for ev, err:=msg.w.ReadEvent(); err==nil; ev, err=msg.w.ReadEvent() {
		if ev.Origin!=goacme.Mouse {
			msg.w.UnreadEvent(ev)
			continue
		}
		quote:=false
		replyall:=false
		if (ev.Type&goacme.Execute)==goacme.Execute {
			switch ev.Text {
				case "Del":
					msg.w.UnreadEvent(ev)
					msg.w.Close()
					msg.w=nil
					return
				case "Delmesg":
					if !msg.deleted {
						msg.deleted=true
						msg.box.deleted++
						msg.w.Del(true)
						msg.w.Close()
						msg.w=nil
						@<Refresh |msg|@>
						return
					}
					continue
				case "UnDelmesg":
					if msg.deleted {
						msg.deleted=false
						msg.box.deleted--
						@<Write a tag of message window@>
						@<Refresh |msg|@>
					}
					continue
				case "Text":
					if len(msg.text)!=0 && len(msg.html)!=0 {
						msg.showhtml=false
						msg.open()
					}
					continue
				case "Html":
					if len(msg.text)!=0 && len(msg.html)!=0 {
						msg.showhtml=true
						msg.open()
					}
					continue
				case "Browser":
					@<Save stuff on disk and plumb a message to a web browser@>
					continue
				case "Save":
					@<Save a message@>
					continue
				case "Q":
					quote=true
					if args:=strings.Fields(ev.Arg); len(args)>0 {
						ev.Text=args[0]
						ev.Arg=strings.Join(args, " ")
					}
					fallthrough
				case "Reply", "Replyall":
					if ev.Text=="Reply" {
						args:=strings.Fields(ev.Arg)
						for _, v:=range args {
							if v=="all" {
								replyall=true
							}
						}
					} else if ev.Text=="Replyall" {
						replyall=true
					}
					@<Compose a message@>
					continue
				case "Up":
					@<Get |parent| for |msg|@>
					if parent!=nil {
						@<Create |msgs|@>
						name:=parent.box.name
						id:=parent.id
						@<Add a |id| message to |msgs|@>
						@<Send |msgs| for viewing@>
					}
					continue
				case "Down":
					@<Get |children| for |msg|@>
					if len(children)!=0 {
						@<Create |msgs|@>
						name:=children[0].box.name
						id:=children[0].id
						@<Add a |id| message to |msgs|@>
						@<Send |msgs| for viewing@>
					}
					continue
				case "Prev":
					@<Get a previous |pmsg|@>
					if pmsg!=nil {
						@<Create |msgs|@>
						name:=pmsg.box.name
						id:=pmsg.id
						@<Add a |id| message to |msgs|@>
						@<Send |msgs| for viewing@>
					}
					continue
				case "Next":
					@<Get a next |nmsg|@>
					if nmsg!=nil {
						@<Create |msgs|@>
						name:=nmsg.box.name
						id:=nmsg.id
						@<Add a |id| message to |msgs|@>
						@<Send |msgs| for viewing@>
					}
					continue
			}
		} else if (ev.Type&goacme.Look)==goacme.Look  {
		}
		msg.w.UnreadEvent(ev)

	}
}()

@
@<Clean |msg.w| window@>=
glog.V(debug).Infof("clean the '%s/%d' message's window\n", msg.box.name, msg.id )
clean(msg.w)

@
@<Imports@>=
"os/exec"

@
@<Types@>=
file struct {
	name 		string
	mimetype	string
	path 		string
}

@ |text| contains a path in \.{mailfs} from the message root to a text variant of the message.

|html| contains a path in \.{mailfs} from the message root to a html variant of the message.

|showmail| is a flag to show the html variant of the message.

|files| contains |*file| with |path|s in \.{mailfs} from the message root to a file is attached in the message,
|mimetype| and |name| of the file.

|cids| contains a map of |"Content-ID"| on |*file|.

@<Rest of |message| members@>=
text		string
html		string
showhtml	bool
files		[]*file
cids		map[string]*file

@
@<Rest of initialization of |message|@>=
cids: make(map[string]*file),

@ If |text| and |html| is empty we should fill them by |bodyPath|, then we read corresponding variant of the message.
In case of the html variant we print |buf| and start a pipe of external programs |"9p"| and |"htmlfmt"| to print the html
message body to the window. Then we fill |buf| with command to obtain contents of |files|.
@<Compose a body of the message@>=
{
	if len(msg.text)==0 && len(msg.html)==0 {
		if err=msg.bodyPath(bfid, ""); err!=nil {
			glog.Errorf("can't ged a body path of '%d': %v\n", msg.id, err)
		}
		glog.V(debug).Infof("paths for bodies of the '%d' message have been found: text-'%s', html-'%s'\n",
							msg.id, msg.text, msg.html)

	}
	if len(msg.text)!=0 && !msg.showhtml {
		glog.V(debug).Infof("using a path for a text body of the '%d' message: '%s'\n", msg.id, msg.text)
		if buf, err=readAll(bfid, msg.text, buf); err!=nil {
			glog.Errorf("can't read '%s': %v\n", msg.text, err)
			return
		}
	} else if len(msg.html)!=0 {
		glog.V(debug).Infof("using a path for a html body of the '%d' message: '%s'\n", msg.id, msg.html)
		msg.w.Write(buf)
		buf=nil
		c1:=exec.Command("9p", "read", fmt.Sprintf("%s/%s/%d/%s", srv, msg.box.name, msg.id, msg.html))
		c2:=exec.Command( "htmlfmt", "-cutf-8")
		c2.Stdout, _=msg.w.File("body")
		c2.Stdin, err=c1.StdoutPipe()
		if err!=nil {
			glog.Errorf("can't get a stdout pipe: %v\n", err)
			return
		}
		if err=c2.Start(); err!=nil {
			glog.Errorf("can't start 'htmlfmt': %v\n", err)
			return
		}
		if err=c1.Run(); err!=nil {
			glog.Errorf("can't run '9p': %v\n", err)
			c2.Wait()
			return
		}
		if err=c2.Wait();err!=nil {
			glog.Errorf("can't wait of completion 'htmlfmt': %v\n", err)
			return
		}
	}
	@<Get |home| enviroment variable@>
	for _, v:= range msg.files {
		buf=append(buf, fmt.Sprintf("\n===> %s (%s)\n", v.path, v.mimetype)...)
		buf=append(buf, fmt.Sprintf("\t9p read %s/%s/%d/%sbody > '%s/%s'\n", srv, msg.box.name, msg.id, v.path, home, v.name)...)
	}
}

@ |bodyPath| recursively looks for parts of the message to determine text and html variants of the message and attached files.
@c
func (msg *message) bodyPath(bfid *client.Fid, path string) error {
	glog.V(debug).Infof("getting a path for a body of the '%d' message\n", msg.id)
	t, err:=readString(bfid, path+"type")
	if err!=nil {
		return err
	}
	switch t {
		case "message/rfc822",
			"text",
			"text/plain",
			"text/richtext",
			"text/tab-separated-values":
			if len(msg.text)==0 {
				msg.text=path+"body"
				glog.V(debug).Infof("a path for a text body of the '%d' message: '%s'\n", msg.id, t)
				return nil
			}
		case "text/html":
			if len(msg.html)==0 {
				msg.html=path+"body"
				glog.V(debug).Infof("a path for a html body of the '%d' message: '%s'\n", msg.id, t)
				return nil
			}
		case "multipart/mixed",
			"multipart/alternative",
			"multipart/related",
			"multipart/signed",
			"multipart/report":
			for c:=1;;c++ {
				if err=msg.bodyPath(bfid, fmt.Sprintf("%s%d/", path, c)); err!=nil {
					break
				}
			}
			return nil
	}
	glog.V(debug).Infof("trying to read '%d/%sfilename'\n", msg.id, path)
	if n, err:=readString(bfid, path+"filename"); err==nil {
		f:=&file{name:n, mimetype:t, path:path,}
		if len(n)==0 {
			f.name="attachment"
		} else if cid, ok:=msg.getCID(path); ok {
			msg.cids[cid]=f
		}
		msg.files=append(msg.files, f)
	}
	return nil
}

@ |getCID| parses |"mimeheader"| and takes |"Content-ID"| identifier for |path|
@c
func (msg *message) getCID(path string) (string, bool) {
	src:=fmt.Sprintf("%d/%smimeheader", msg.id, path)
	glog.V(debug).Infof("getting of cids for path '%s'\n", src)
	fid, err:=msg.box.fid.Walk(src)
	if err==nil {
		err=fid.Open(plan9.OREAD)
	}
	if err!=nil {
		glog.Errorf("can't open '%s': %v\n", src, err)
		return "", false
	}
	defer fid.Close()
	fid.Seek(0, 0)
	b:=bufio.NewReader(fid)
	for s, err:=b.ReadString('\n'); err==nil||err==io.EOF; s, err=b.ReadString('\n') {
		glog.V(debug).Infof("looking for a cid in '%s'\n", s)
		if strings.HasPrefix(s, "Content-ID: <") {
			s=s[len("Content-ID: <"):len(s)-2]
			glog.V(debug).Infof("found a cid '%s'\n", s)
			return s, true
		}
		if err==io.EOF {
			break
		}
	}
	return "", false
}



@ |home| environment variable.
@<Variables@>=
home	string

@
@<Get |home| enviroment variable@>=
@<Get some things at once@>

@
@<Get it at once@>=
if home=os.Getenv("home"); len(home)==0 {
	if home=os.Getenv("HOME"); len(home)==0 {
		glog.Errorln("can't get a home directory from the environment, the home is assumed '/'")
		home="/"
	}
}


@ |readStrings| reads a full string from |name| file with |pfid| like a root.
@c
func readString(pfid *client.Fid, name string) (str string, err error) {
	glog.V(debug).Infof("readString: trying to open '%s'\n", name)
	f, err:=pfid.Walk(name)
	if err==nil {
		err=f.Open(plan9.OREAD)
	}
	if err!=nil {
		return
	}
	defer f.Close()
	str, err=bufio.NewReader(f).ReadString('\n')
	if err!=nil && err!=io.EOF {
		return
	}
	return str, nil
}

@ |readAll| reads all content of |name| file with |pfid| like a root in |buf|
@c
func readAll(pfid *client.Fid, name string, buf []byte) ([]byte, error) {
	glog.V(debug).Infof("readAll: trying to open '%s'\n", name)
	f, err:=pfid.Walk(name)
	if err==nil {
		err=f.Open(plan9.OREAD)
	}
	if err!=nil {
		return buf, err
	}
	defer f.Close()
	b:=bufio.NewReader(f)
	for s, err:=b.ReadString('\n'); err==nil||err==io.EOF; s, err=b.ReadString('\n') {
		if strings.HasSuffix(s, "\r\n") {
			s=strings.TrimRight(s, "\r\n")
			s+="\n"
		}
		buf=append(buf, s...)
		if err==io.EOF {
			break
		}
	}
	return buf, nil
}


@ To view a message in a web brower we need to store a body of the message and all images of the message on disk
and plumb a full pathname of the message to |"web"| rule. But in case of the images the body should be fixed
to help a browser to find the images.

@<Save stuff on disk and plumb a message to a web browser@>=
{
	@<Get current |user|@>
	dir:=fmt.Sprintf("%s/amail-%s/%s/%d", os.TempDir(), cuser, msg.box.name, msg.id)
	if err:=os.MkdirAll(dir, 0700); err!=nil {
		glog.Errorf("can't create a directory '%s': %v\n", dir, err)
		continue
	}

	if len(msg.files)==0 {
		if err:=saveFile(fmt.Sprintf("%s/%s/%d/%s", srv, msg.box.name, msg.id, msg.html),
						fmt.Sprintf("%s/%d.html", dir, msg.id)); err !=nil {
			continue
		}
	} else {
		if err:=msg.fixFile(dir); err !=nil {
			continue
		}
		for _, v:=range msg.files {
 			saveFile(fmt.Sprintf("%s/%s/%d/%s/body", srv, msg.box.name, msg.id, v.path),
 						fmt.Sprintf("%s/%s", dir, v.name))
		}

	}

	if p, err:=goplumb.Open("send", plan9.OWRITE); err!=nil {
		glog.Errorf("can't open plumbing port 'send': %v\n", err)
	} else if err:=p.SendText("amail", "web", dir, fmt.Sprintf("file://%s/%d.html", dir, msg.id)); err!=nil {
		glog.Errorf("can't plumb a message '%s': %v\n", fmt.Sprintf("file://%s/%d.html", dir, msg.id), err)
	}
}

@ |saveFile| saves file on a disk by call |"9p"|
@c
func saveFile(src, dst string) error {
	var err error
	c:=exec.Command("9p", "read", src)
	f, err:=os.OpenFile(dst, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err!=nil {
		glog.Errorf("can't create a file '%s': %v\n", dst, err)
		return err
	}
	defer f.Close()
	c.Stdout=f
	if err=c.Run(); err!=nil {
		glog.Errorf("can't run '9p': %v\n", err)
	}
	return err
}

@ |fixFile| reads the message body  and replaces all |"cid"| on corresponding cids.
@c
func (msg *message) fixFile(dir string) error {
	src:=fmt.Sprintf("%d/%s", msg.id, msg.html)
	dst:=fmt.Sprintf("%s/%d.html", dir, msg.id)
	df, err:=os.OpenFile(dst, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err!=nil {
		glog.Errorf("can't create a file '%s': %v\n", dst, err)
		return err
	}
	defer df.Close()
	fid, err:=msg.box.fid.Walk(src)
	if err==nil {
		err=fid.Open(plan9.OREAD)
	}
	if err!=nil {
		glog.Errorf("can't open to '%s': %v\n", src, err)
		return err
	}
	defer fid.Close()
	b:=bufio.NewReader(fid)
	for s, err:=b.ReadString('\n'); err==nil||err==io.EOF; s, err=b.ReadString('\n') {
		p:=0
		for b:=strings.Index(s[p:], "\"cid:"); b!=-1; b=strings.Index(s[p:], "\"cid:") {
			b+=p
			e:=strings.Index(s[b+1:], "\"")
			if e==-1 {
				break
			}
			e++
			glog.V(debug).Infof("len(s): %v, p: %v, b: %v, e: %v\n", len(s), p, b, e)
			cid:=s[b+5:b+e]
			glog.V(debug).Infof("cid: %s\n", cid)
			if f, ok:=msg.cids[cid]; ok {
				glog.V(debug).Infof("found a cid: %s, replace '%s' by '%s'\n", cid, s[b+1:b+e], f.name )
				s=strings.Replace(s, s[b+1:b+e], f.name, 1)
			} else {
				p=b+e
			}
		}
		df.Write([]byte(s))
		if err==io.EOF {
			break
		}
	}
	return err
}

@
@<Imports@>=
"os/user"

@ current |user|
@<Variables@>=
cuser	string

@
@<Get current |user|@>=
@<Get some things at once@>

@
@<Get it at once@>=
if u, err:=user.Current(); err!=nil {
	glog.Errorf("can't get a name of the current user: %v\n", err)
} else {
	cuser=u.Username
}


@
@<Imports@>=
"sync"

@
@<Variables@>=
once sync.Once

@
@<Get some things at once@>=
once.Do(func() {@<Get it at once@>})

@
@<Save a message@>=
{
	if len(ev.Arg)==0 {
		continue
	}
	f, err:=msg.box.fid.Walk("ctl")
	if err==nil {
		err=f.Open(plan9.OWRITE)
	}
	if err!=nil {
		glog.Errorf("can't open 'ctl': %v\n", err)
		continue
	}
	bs:=strings.Fields(ev.Arg)
	for _, v:=range bs {
		s:=fmt.Sprintf("save %s %d/", v, msg.id)
		if _, err:=f.Write([]byte(s)); err!=nil {
			glog.Errorf("can't write '%s' to 'ctl': %v\n", s, err)
		}
	}
	f.Close()
}


@* Composing a message.
@<Compose a message@>=
{
	@<Create a new message window@>
	name:=fmt.Sprintf("Amail/%s/%d/%sReply%s", @t\1@>@/
							msg.box.name, @/
							msg.id, @/
							func()string{if quote {return "Q"}; return ""}(), @/
							func()string{if replyall {return "all"}; return ""}()@t\2@>)
	@<Print the |name| for window |w|@>
	buf:=make([]byte, 0, 0x8000)
	buf=append(buf, fmt.Sprintf("To: %s\n", msg.from)...)
	if replyall {
		for _, v:=range msg.to {
			buf=append(buf, fmt.Sprintf("To: %s\n", v)...)
		}
		for _, v:=range msg.cc {
			buf=append(buf, fmt.Sprintf("To: %s\n", v)...)
		}
	}
	buf=append(buf, fmt.Sprintf("Subject: %s%s\n", @t\1@>@/
		func() string{
			if !strings.Contains(msg.subject, "Re:") {
				return "Re: "
			}
			return ""
		}(), @/
		msg.subject)...@t\2@>)
	if quote {
		buf=append(buf, '\n')
		@<Add quoted message@>
	} else {
		buf=append(buf, fmt.Sprintf("Include: Mail/%s/%d/raw\n", msg.box.name, msg.id)...)
		@^Using of \.{Mail} is required by \.{upas/marshal}@>
	}
	buf=append(buf, '\n')
	w.Write(buf)
	@<Append |msg.box|-specific signature@>
}

@
@<Create a new message window@>=
w, err:=goacme.New()
if err!=nil {
	glog.Errorf("can't create a window: %v\n", err)
	continue
}
if err:=writeTag(w, " Look Post Undo "); err!=nil {
	glog.Errorf("can't write a tag for a new message window: %v\n", err)
}
@<Start a goroutine to process events from a composed mail window@>

@ If we are going to reply a message, we should specify |msg|.
@<Start a goroutine to process events from a composed mail window@>=
go func(msg *message) {
	glog.V(debug).Infoln("starting a goroutine to process events from a composed mail window")
	for ev, err:=w.ReadEvent(); err==nil; ev, err=w.ReadEvent() {
		if ev.Origin!=goacme.Mouse {
			w.UnreadEvent(ev)
			continue
		}
		if (ev.Type&goacme.Execute)==goacme.Execute {
			switch ev.Text {
				case "Del":
					w.UnreadEvent(ev)
					w.Close()
					return
				case "Post":
					@<Send the message@>
			}
		}
		w.UnreadEvent(ev)
	}
}(msg)

@
@<Add quoted message@>=
if len(msg.text)!=0 {
	fn:=fmt.Sprintf("%d/%s", msg.id, msg.text)
	f, err:=msg.box.fid.Walk(fn)
	if err==nil {
		err=f.Open(plan9.OREAD)
	}
	if err!=nil {
		glog.Errorf("can't open '%s/%s/%s': %v\n", srv, msg.box.name, fn)
		continue
	}
	@<Quote a message@>
	f.Close()
} else if len(msg.html)!=0 {
	@<Quote a html message@>
}

@ To quote a message we read strings from |f| and add |"> "| to the begin of every string.
@<Quote a message@>=
{
	b:=bufio.NewReader(f)
	for s, err:=b.ReadString('\n'); err==nil || err==io.EOF; s, err=b.ReadString('\n') {
		buf=append(buf, '>', ' ')
		if strings.HasSuffix(s, "\r\n") {
			s=strings.TrimRight(s, "\r\n")
			s+="\n"
		}
		buf=append(buf, s...)
		if err==io.EOF {
			break
		}
	}
}

@ To quote the html message we start a pipe of external programs |"9p"| and |"htmlfmt"| and read an output of |"htmlfmt"|
@<Quote a html message@>=
{
	c1:=exec.Command("9p", "read", fmt.Sprintf("%s/%s/%d/%s", srv, msg.box.name, msg.id, msg.html))
	c2:=exec.Command( "htmlfmt", "-cutf-8")
	f, err:=c2.StdoutPipe()
	if err!=nil {
		glog.Errorf("can't get a stdout pipe: %v\n", err)
	}
	c2.Stdin, err=c1.StdoutPipe()
	if err!=nil {
		glog.Errorf("can't get a stdout pipe: %v\n", err)
		f.(io.Closer).Close()
		continue
	}
	if err=c2.Start(); err!=nil {
		glog.Errorf("can't start 'htmlfmt': %v\n", err)
		f.(io.Closer).Close()
		continue
	}
	if err=c1.Start(); err!=nil {
		glog.Errorf("can't run '9p': %v\n", err)
		c2.Wait()
		f.(io.Closer).Close()
		continue
	}
	@<Quote a message@>
	c1.Wait()
	c2.Wait()
	f.(io.Closer).Close()
}

@ To send a message we start an external program |"upas/marshal"| and send to its input recipient and body of the message.
If |msg!=nil|, it will be added like a message is replied.
@<Send the message@>=
{
	@<Get |plan9dir| from enviroment variable@>
	w.Seek(0, 0)
	w.WriteAddr(wholefile)
	ff, _:=w.File("xdata")
	b:=bufio.NewReader(ff)
	var to, cc, bcc, attach, include []string
	var subject string
	for {
		s, err:=b.ReadString('\n')
		if err!=nil {
			break
		}
		s=strings.TrimSpace(s)
		if len(s)==0 {
			// an empty line, the rest is a body of the message
			break
		}
		p:=strings.Index(s, ":")
		if p!=-1 {
			f:=strings.Split(s[p+1:], ",")
			for i, _:=range f {
				f[i]=strings.TrimSpace(f[i])
			}
			switch strings.ToLower(s[:p]) {
				case "to":
					@<Get last elements of addresses@>
					to=append(to, f...)
				case "cc":
					@<Get last elements of addresses@>
					cc=append(cc, f...)
				case "bcc":
					@<Get last elements of addresses@>
					bcc=append(bcc, f...)
				case "attach":
					attach=append(attach, f...)
				case "include":
					include=append(include, f...)
				case "subject":
					subject=fmt.Sprintf("%s", strings.TrimSpace(s[p+1:]))
			}
		} else {
			// recipient addresses can be written without "to:"
			f:=strings.Split(s, ",")
			for i, _:=range f {
				f[i]=strings.TrimSpace(f[i])
			}
			@<Get last elements of addresses@>
			to=append(to, f...)
		}
	}
	args:=append([]string{}, "-8")
	if msg!=nil{
		args=append(args, "-R", fmt.Sprintf("%s/%d", msg.box.name, msg.id))
	}
	if len(subject)!=0 {
		args=append(args, "-s", subject)
	}
	for _, v:=range include {
		args=append(args, "-A", v)
	}
	for _, v:=range attach {
		args=append(args, "-a", v)
	}
	c:=exec.Command(plan9dir+"/bin/upas/marshal", args...)
	p, err:=c.StdinPipe()
	if err!=nil {
		glog.Errorf("can't get a stdin pipe: %v\n", err)
		continue
	}
	if err:=c.Start(); err!=nil {
		glog.Errorf("can't start 'upas/marshal': %v\n", err)
		continue
	}
	if len(to)!=0 {
		if _, err:=fmt.Fprintln(p, "To:", strings.Join(to, ",")); err!=nil {
			glog.Errorf("can't write 'to' fields to 'upas/marshal': %v\n", err)
			continue
		}
	}
	glog.V(debug).Infoln("to is written")
	if len(cc)!=0 {
		if _, err:=fmt.Fprintln(p, "CC:", strings.Join(cc, ",")); err!=nil {
			glog.Errorf("can't write 'cc' fields to 'upas/marshal': %v\n", err)
			continue
		}
	}
	glog.V(debug).Infoln("cc is written")
	if len(bcc)!=0 {
		if _, err:=fmt.Fprintln(p, "BCC:", strings.Join(bcc, ",")); err!=nil {
			glog.Errorf("can't write 'bcc' fields to  'upas/marshal': %v\n", err)
			continue
		}
	}
	glog.V(debug).Infoln("bcc is written")
	for s, err:=b.ReadString('\n'); err==nil || err==io.EOF; s, err=b.ReadString('\n') {
		glog.V(debug).Infof("writing '%s':%v", s, err)

		p.Write([]byte(s))
		if err==io.EOF {
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

@ An address can be preffered by name or alias of recipient. It has to be removed.
@<Get last elements of addresses@>=
for i, _:=range f {
	f[i]=strings.TrimSpace(f[i])
	if sf:=strings.Fields(f[i]); len(sf)>1 {
		f[i]=strings.TrimSpace(sf[len(sf)-1])
	}
}


@
@<Variables@>=
plan9dir string

@
@<Get |plan9dir| from enviroment variable@>=
@<Get some things at once@>

@
@<Get it at once@>=
if plan9dir=os.Getenv("PLAN9"); len(plan9dir)==0 {
	glog.Errorln("can't get PLAN9 directory from the environment, the plan9dir is assumed '/usr/local/plan9'")
	plan9dir="/usr/local/plan9"
}

@
@<Append |box|-specific signature@>=
writeSignature(w, box)

@
@<Append |msg.box|-specific signature@>=
if msg!=nil {
	writeSignature(w, msg.box)
} else {
	writeSignature(w, nil)
}

@
@<Append common signature@>=
writeSignature(w, nil)


@ At first we are looking for |box|-specific signature in |$HOME/mail/<mailbox>.signature| file.
If the file doesn't exist, we are trying to open |$HOME/mail/<mailbox>.signature| file with common signature.
@c
func writeSignature(w *goacme.Window, box *mailbox) {
	@<Get |home| enviroment variable@>
	var f io.ReadCloser
	var err error
	if box!=nil {
		f,err=os.Open(fmt.Sprintf("%s/mail/%s.signature", home, box.name))
	}
	if err!=nil || f==nil {
		f,err=os.Open(fmt.Sprintf("%s/mail/signature", home))
	}
	if err==nil {
		w.Write([]byte("\n"))
		b:=bufio.NewReader(f)
		for buf, err:=b.ReadBytes('\n'); err==nil || err==io.EOF; buf, err=b.ReadBytes('\n') {
			w.Write(buf)
			if err==io.EOF {
				break
			}
		}
		f.Close()
	}
	@<Go to top of window |w|@>
}

@** Index.
