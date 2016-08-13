FROM microsoft/dotnet
ENV GOLANG_VERSION 1.7rc6
ENV GOLANG_DOWNLOAD_URL https://golang.org/dl/go$GOLANG_VERSION.linux-amd64.tar.gz
ENV GOLANG_DOWNLOAD_SHA256 45e3dfba542927ea58146a5d47a983feb36401ccafeea28a9e0a79534738b154
ENV GOPATH /go
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH

ADD root/ /root/

RUN apt-get update && apt-get install -y --no-install-recommends \
			g++ \
			gcc \
			libc6-dev \
			make \
			ncurses-dev \
			libtolua-dev \
			exuberant-ctags \
	&& rm -rf /var/lib/apt/lists/*                                          \
  && curl -fsSL "$GOLANG_DOWNLOAD_URL" -o golang.tar.gz                   \
	&& echo "$GOLANG_DOWNLOAD_SHA256  golang.tar.gz" | sha256sum -c -       \
	&& tar -C /usr/local -xzf golang.tar.gz                                 \
	&& rm golang.tar.gz                                                     \
	&& mkdir -p "$GOPATH/src" "$GOPATH/bin" && chmod -R 777 "$GOPATH"       \
  && ln -s /usr/include/lua5.2/ /usr/include/lua                          \
	&& ln -s /usr/lib/x86_64-linux-gnu/liblua5.2.so /usr/lib/liblua.so      \
	&& cd /tmp                                                              \
  && git clone https://github.com/vim/vim.git                             \
  && cd vim                                                               \
  && ./configure --with-features=huge --enable-luainterp                  \
        --enable-gui=no --without-x --prefix=/usr                         \
  && make VIMRUNTIMEDIR=/usr/share/vim/vim74                              \
  && make install                                                         \
	&& mkdir -p ~/.vim/bundle                                               \
  && cd  ~/.vim/bundle                                                    \
  && git clone --depth 1 https://github.com/gmarik/Vundle.vim.git         \
  && git clone --depth 1 https://github.com/fatih/vim-go.git              \
  && git clone --depth 1 https://github.com/majutsushi/tagbar.git         \
  && git clone --depth 1 https://github.com/Shougo/neocomplete.vim.git    \
  && git clone --depth 1 https://github.com/scrooloose/nerdtree.git       \
  && git clone --depth 1 https://github.com/bling/vim-airline.git         \
  && git clone --depth 1 https://github.com/tpope/vim-fugitive.git        \
  && git clone --depth 1 https://github.com/jistr/vim-nerdtree-tabs.git   \
  && git clone --depth 1 https://github.com/mbbill/undotree.git           \
  && git clone --depth 1 https://github.com/Lokaltog/vim-easymotion.git   \
  && git clone --depth 1 https://github.com/scrooloose/nerdcommenter.git  \
	&& rm -rf Vundle.vim/.git vim-go/.git tagbar/.git neocomplete.vim/.git  \
        nerdtree/.git vim-airline/.git vim-fugitive/.git                  \
        vim-nerdtree-tabs/.git undotree/.git vim-easymotion/.git          \
        nerdcommenter/.git

RUN vim +PluginInstall +qall
RUN vim +GoInstallBinaries +qall

# RUN apt-key adv --keyserver hkp://keyserver.ubuntu.com:80 --recv-keys 3FA7E0328081BFF6A14DA29AA6A19B38D3D831EF \
#  && echo "deb http://download.mono-project.com/repo/debian wheezy/snapshots/4.4.2.11 main" > /etc/apt/sources.list.d/mono-xamarin.list \
#	&& echo "deb http://download.mono-project.com/repo/debian wheezy-apache24-compat main" | tee -a /etc/apt/sources.list.d/mono-xamarin.list \
#  && echo "deb http://download.mono-project.com/repo/debian wheezy-libjpeg62-compat main" | tee -a /etc/apt/sources.list.d/mono-xamarin.list \
#  && apt-get update \
#  && apt-get install --fix-missing -y mono-complete \
#  && rm -rf /var/lib/apt/lists/* /tmp/*
