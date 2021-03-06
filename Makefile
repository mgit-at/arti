##
## Copyright: 2016, mgIT GmbH <office@mgit.at>
## All rights reserved.
##
## Licensed under the Apache License, Version 2.0 (the "License");
## you may not use this file except in compliance with the License.
## You may obtain a copy of the License at
##
## http:##www.apache.org/licenses/LICENSE-2.0
##
## Unless required by applicable law or agreed to in writing, software
## distributed under the License is distributed on an "AS IS" BASIS,
## WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
## See the License for the specific language governing permissions and
## limitations under the License.
##

bin:
	go build -tags netgo

clean:
	rm -f arti
	rm -f manpage/manpage
	rm -f manpage/*.3

update:
	glide update -u -s -v --cache

manpage:
	@cd manpage/ &&	go build && ./manpage

deb: bin manpage
	eatmydata debuild -I -us -uc

.PHONY: bin clean update manpage deb
