# Install quarto
apt-get update
curl -o quarto.deb -L "https://github.com/quarto-dev/quarto-cli/releases/download/v1.4.554/quarto-1.4.554-linux-amd64.deb"
dpkg -i quarto.deb 
rm quarto.deb

# Install quarto dependencies
quarto install tinytex                                  # TinyTex
quarto install chromium                                 # Chromium

apt-get install -y python3-pip                          # Pip
python3 -m pip install jupyter --break-system-packages  # Jupyter

apt install -y r-base r-base-dev                        # R
Rscript -e 'install.packages("knitr")'                  # Knitr
Rscript -e 'install.packages("rmarkdown")'              # Rmarkdown
