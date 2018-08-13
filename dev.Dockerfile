FROM golang:1.10.3-stretch AS qsub

# Build the qsub substitute.
RUN git clone https://github.com/alejandrox1/falcon "$GOPATH/src/github.com/alejandrox1/falcon" \
    && go build -o /usr/bin/qsub github.com/alejandrox1/falcon \
    && rm -rf $GOPATH


FROM ubuntu:18.04
# Falcon :-( 
# https://github.com/PacificBiosciences/FALCON-integrate/wiki/Installation

WORKDIR /src

# Download miniconda, samtools, falcon_unzip, and mummer releases.
ADD  https://repo.continuum.io/miniconda/Miniconda3-latest-Linux-x86_64.sh .
ADD https://github.com/samtools/samtools/releases/download/1.9/samtools-1.9.tar.bz2 .
ADD https://downloads.pacbcloud.com/public/falcon/falcon-2018.03.12-04.00-py2.7-ucs4.tar.gz .
ADD https://github.com/mummer4/mummer/releases/download/v4.0.0beta2/mummer-4.0.0beta2.tar.gz .

ENV FALCON_PREFIX="/fc_env"
ENV PATH="/opt/conda/bin:$PATH"

# Install system dependencies.
RUN apt-get update -y && apt-get install -y git gcc g++ make bzip2 wget \
    libz-dev libncurses5-dev libncursesw5-dev libbz2-dev liblzma-dev \
    libhdf5-serial-dev \
    net-tools \
    && bash Miniconda3-latest-Linux-x86_64.sh -b -p /opt/conda \
    && conda install -y python=2.7.9 \
    && conda install -y virtualenv \
    && mkdir $FALCON_PREFIX 

# Install samtools and other falcon dependencies.
RUN tar xvjf samtools-1.9.tar.bz2 \                                          
    && cd samtools-1.9 \                                                        
    && ./configure \                                                            
    && make \                                                                   
    && make install \                                                           
    && cd /src \                                                                 
    && git clone https://github.com/lh3/minimap2 \                              
    && cd minimap2 && make \                                                    
    && conda install -y -c bioconda blasr \                                     
    && conda install -y numpy h5py \                                            
    && pip install git+https://github.com/PacificBiosciences/pbalign

# Install falcon.
RUN virtualenv $FALCON_PREFIX \
    && tar xvzf falcon-2018.03.12-04.00-py2.7-ucs4.tar.gz -C $FALCON_PREFIX \
    && tar zxvf mummer-4.0.0beta2.tar.gz \
    && cd mummer-4.0.0beta2 \
    && ./configure \
    && make && make install \
    && cd src/ \
    && git clone git://github.com/PacificBiosciences/FALCON-integrate.git \
    && . $FALCON_PREFIX/bin/activate

ENV LD_LIBRARY_PATH="${FALCON_PREFIX}/lib:${LD_LIBRARY_PATH}"
ENV PATH="${FALCON_PREFIX}/bin:/minimap2:${PATH}"

# Add the qsub substitute into this image.
COPY --from=qsub /usr/bin/qsub /usr/bin/


# Falcon (interactive) tests: fc_run fc_run_ecoli_2.cfg
WORKDIR /ecoli_test
ADD https://raw.githubusercontent.com/PacificBiosciences/FALCON/master/examples/fc_run_ecoli_2.cfg .
WORKDIR /ecoli_test/data
ADD https://www.dropbox.com/s/tb78i5i3nrvm6rg/m140913_050931_42139_c100713652400000001823152404301535_s1_p0.1.subreads.fasta .
ADD https://www.dropbox.com/s/v6wwpn40gedj470/m140913_050931_42139_c100713652400000001823152404301535_s1_p0.2.subreads.fasta .
ADD https://www.dropbox.com/s/j61j2cvdxn4dx4g/m140913_050931_42139_c100713652400000001823152404301535_s1_p0.3.subreads.fasta .
WORKDIR /ecoli_test


CMD ["/bin/bash"]
