
import sys
import bz2
import xml.sax.handler

import samoa



class WikiHandler(xml.sax.handler.ContentHandler):

    def __init__(self):
        xml.sax.handler.ContentHandler.__init__(self)
        self.in_title = False

    def startElement(self, name, attrs):
        if name == 'title':
            self.in_title = True
        #print "startElement(%r, %r)" % (name, attrs)

    def endElement(self, name):
        self.in_title = False
        #print "endElement(%r)" % name

    def startElementNS(self, name, attrs):
        #print "startElementNS(%r, %r)" % (name, attrs)
        pass

    def endElementNS(self, name):
        #print "endElementNS(%r)" % name
        pass

    def characters(self, content):
        if self.in_title:
            print content
        #print "content(%r)" % content[:120]
        pass

    def ignorableWhitespace(self, content):
        #print "ignorableWhitespace(%r)" % content[:120]
        pass

    def processingInstruction(self, target, data):
        #print "processingInstruction(%r, %r)" % (target, data)
        pass

    def skippedEntity(self, name):
        #print "skippedEntity(%r)" % name
        pass


dump_input_path = sys.argv[1]
dump_in = bz2.BZ2File(dump_input_path)

xml_parser = xml.sax.make_parser()
xml_parser.setContentHandler(WikiHandler())

xml_parser.parse(dump_in)

