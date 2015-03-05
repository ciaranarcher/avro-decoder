require 'avro'

file = File.open('data.avro', 'r+')
dr = Avro::DataFile::Reader.new(file, Avro::IO::DatumReader.new)
dr.each { |record| p record }
