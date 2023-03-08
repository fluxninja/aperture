package com.javademoapp.javademoapp;

import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.fasterxml.jackson.databind.module.SimpleModule;
import com.fasterxml.jackson.databind.ser.std.StdSerializer;

import java.io.IOException;
import java.util.ArrayList;
import java.util.List;

public class SubrequestChain {
    private List<Subrequest> subrequests = new ArrayList<>();

    public List<Subrequest> getSubrequests() {
        return subrequests;
    }

    public void setSubrequests(List<Subrequest> subrequests) {
        this.subrequests = subrequests;
    }

    public void addSubrequest(Subrequest subrequest) {
        this.subrequests.add(subrequest);
    }

    // Custom serializer to write subrequests as a JSON list
    public static class SubrequestChainSerializer extends StdSerializer<SubrequestChain> {
        public SubrequestChainSerializer() {
            this(null);
        }

        public SubrequestChainSerializer(Class<SubrequestChain> t) {
            super(t);
        }

        @Override
        public void serialize(SubrequestChain value, com.fasterxml.jackson.core.JsonGenerator gen, com.fasterxml.jackson.databind.SerializerProvider provider) throws IOException {
            gen.writeStartArray();
            for (Subrequest subrequest : value.getSubrequests()) {
                gen.writeObject(subrequest);
            }
            gen.writeEndArray();
        }
    }

    // Custom deserializer to read subrequests as a JSON list
    public static class SubrequestChainDeserializer extends com.fasterxml.jackson.databind.JsonDeserializer<SubrequestChain> {
        @Override
        public SubrequestChain deserialize(com.fasterxml.jackson.core.JsonParser p, com.fasterxml.jackson.databind.DeserializationContext ctxt) throws IOException, JsonProcessingException {
            List<Subrequest> subrequests = new ArrayList<>();
            while (p.nextToken() != com.fasterxml.jackson.core.JsonToken.END_ARRAY) {
                Subrequest subrequest = p.readValueAs(Subrequest.class);
                subrequests.add(subrequest);
            }
            SubrequestChain subrequestChain = new SubrequestChain();
            subrequestChain.setSubrequests(subrequests);
            return subrequestChain;
        }
    }

    // Register custom serializer and deserializer with ObjectMapper
    static ObjectMapper mapper = new ObjectMapper();
    static SimpleModule module = new SimpleModule();
    static {
        module.addSerializer(SubrequestChain.class, new SubrequestChainSerializer());
        module.addDeserializer(SubrequestChain.class, new SubrequestChainDeserializer());
        mapper.registerModule(module);
    }

    // Convert to JSON string
    public String toJson() throws JsonProcessingException {
        return mapper.writeValueAsString(this);
    }

    // Parse from JSON string
    public static SubrequestChain fromJson(String json) throws IOException {
        return mapper.readValue(json, SubrequestChain.class);
    }
}
